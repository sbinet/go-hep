package io

import (
	"archive/tar"
	"compress/gzip"
	"encoding/gob"
	"errors"
	"io"
	//"json"
	"os"
	"path/filepath"
	"strings"
)

const (
	_HEPIO_ROOT = "__go-hep-io__"
	_HEADER     = "__header__"
)

type adict map[string]interface{}

type writerz struct {
	enc   *gob.Encoder
	compr *gzip.Writer
}

type readerz struct {
	dec   *gob.Decoder
	compr *gzip.Reader
}

type File struct {
	dir    string
	hdr    *os.File
	tuples map[string]*Tuple
}

func Create(fname string) (*File, error) {
	dir_root := filepath.Join(_HEPIO_ROOT, fname)
	err := os.RemoveAll(dir_root)
	if err != nil {
		return nil, err
	}
	dirname, err := filepath.Abs(dir_root)
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(dirname, 0776)
	if err != nil {
		return nil, err
	}
	// create header file
	f, err := os.Create(filepath.Join(dirname, _HEADER))
	if err != nil {
		return nil, err
	}
	hepfile := &File{
		dir:    dirname,
		hdr:    f,
		tuples: make(map[string]*Tuple),
	}
	return hepfile, nil
}

func Open(fname string) (*File, error) {
	dir_root := filepath.Join(_HEPIO_ROOT, fname)
	err := os.RemoveAll(dir_root)
	if err != nil {
		return nil, err
	}
	dirname, err := filepath.Abs(dir_root)
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(dirname, 0776)
	if err != nil {
		return nil, err
	}

	// open hepfile
	raw_f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	tr := tar.NewReader(raw_f)

	// create local header file
	f, err := os.Create(filepath.Join(dirname, _HEADER))
	if err != nil {
		return nil, err
	}
	// untar everything
	// FIXME: make this lazy...
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}
		if err != nil {
			return nil, err
		}
		if hdr.Name == _HEADER {
			_, err = io.Copy(f, tr)
			if err != nil {
				return nil, err
			}
		} else {
			ff, err := os.Create(filepath.Join(dirname,
				filepath.Base(hdr.Name)))
			if err != nil {
				return nil, err
			}
			_, err = io.Copy(ff, tr)
			if err != nil {
				return nil, err
			}
			err = ff.Close()
			if err != nil {
				return nil, err
			}
		}
	}
	hepfile := &File{
		dir:    dirname,
		hdr:    f,
		tuples: make(map[string]*Tuple),
	}
	return hepfile, nil
}

// Create a new Tuple in folder name 
func (f *File) CreateTuple(name string) (*Tuple, error) {
	ff, err := os.Create(filepath.Join(f.dir, name))
	if err != nil {
		return nil, err
	}
	w, err := gzip.NewWriterLevel(ff, gzip.DefaultCompression)
	if err != nil {
		return nil, err
	}
	enc := gob.NewEncoder(w)
	wz := &writerz{enc: enc, compr: w}
	f.tuples[name] = &Tuple{f: ff, w: wz, r: nil, nentries: 0}
	return f.tuples[name], nil
}

func (f *File) OpenTuple(name string) (*Tuple, error) {
	ff, err := os.Open(filepath.Join(f.dir, name))
	if err != nil {
		return nil, err
	}
	//fi,err := ff.Stat()
	//println("file:",ff.Name(),"size:",fi.Size)
	r, err := gzip.NewReader(ff)
	if err != nil {
		return nil, err
	}
	dec := gob.NewDecoder(r)
	rz := &readerz{dec: dec, compr: r}
	// retrieve metadata for this n-tuple
	metadata := f.fetch_metadata_for(name)
	var nentries int64 = 0
	nentries = metadata.Nentries
	f.tuples[name] = &Tuple{f: ff, w: nil, r: rz, nentries: nentries}
	return f.tuples[name], nil
}

func (f *File) fetch_metadata_for(name string) metadata {
	// save current position
	cur, err := f.hdr.Seek(0, 1)
	defer f.hdr.Seek(cur, 0)

	_, err = f.hdr.Seek(0, 0)
	if err != nil {
		panic(err)
	}
	r, err := gzip.NewReader(f.hdr)
	if err != nil {
		panic(err)
	}
	dec := gob.NewDecoder(r)
	// read hepfile metadata
	file_mdata := make(map[string]interface{})
	err = dec.Decode(&file_mdata)
	if err != nil {
		panic(err)
	}
	//println("file metadata:",file_mdata["version"].(uint32))
	mdata := make(map[string]metadata)
	err = dec.Decode(&mdata)
	if err != nil {
		panic(err)
	}
	v, ok := mdata[name]
	if !ok {
		panic("no metadata for tuple [" + name + "]")
	}
	return v
}

func (f *File) Close() error {
	metadata := make(map[string]metadata)
	for n, t := range f.tuples {
		//println("-- closing:",t.f.Name())
		err := t.Close()
		if err != nil {
			return errors.New("problem closing [" + n + "]: " + err.Error())
		}
		metadata[n] = make_metadata_from(t)
	}

	w, err := gzip.NewWriterLevel(f.hdr, gzip.DefaultCompression)
	if err != nil {
		return err
	}
	enc := gob.NewEncoder(w)
	mdata := map[string]interface{}{
		"version": uint32(0x00000001),
	}
	err = enc.Encode(&mdata)
	if err != nil {
		return err
	}
	err = enc.Encode(&metadata)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}

	// now create the final tar-file
	outname := filepath.Base(f.dir)
	out, err := os.Create(outname)
	if err != nil {
		return err
	}
	tw := tar.NewWriter(out)
	// write top directory
	s := strings.Split(outname, ".")
	root_name := strings.Join(s[:len(s)-1], ".")
	fi, err := out.Stat()
	if err != nil {
		println(err)
		return err
	}
	hdr := &tar.Header{
		Name: root_name,
		Mode: int64(492),
		//Uid:      fi.Uid,
		//Gid:      fi.Gid,
		Size:     0,
		Typeflag: tar.TypeDir,
		//Mtime:    fi.Mtime_ns / 1e9, // ns to s
		//Atime:    fi.Atime_ns / 1e9,
		//Ctime:    fi.Ctime_ns / 1e9,
	}

	err = tw.WriteHeader(hdr)
	if err != nil {
		return err
	}
	tarfile := func(fd *os.File) error {
		fd.Sync()
		tuple, err := os.Open(fd.Name())
		if err != nil {
			println(err.Error())
			return err
		}
		stats, err := tuple.Stat()
		if err != nil {
			println(err)
			return err
		}
		hdr := &tar.Header{
			Name: filepath.Join(root_name, filepath.Base(fd.Name())),
			Mode: int64(fi.Mode()),
			//Uid:      stats.Uid,
			//Gid:      stats.Gid,
			Size:     stats.Size(),
			Typeflag: tar.TypeReg,
			ModTime:  stats.ModTime(),
			//Atime:    stats.Atime_ns / 1e9, // to seconds
			//Ctime:    stats.Ctime_ns / 1e9, // to seconds
		}
		//println("-->",hdr.Name,hdr.Size)
		err = tw.WriteHeader(hdr)
		if err != nil {
			println(err)
			return err
		}
		io.Copy(tw, tuple)
		return nil
	}
	err = f.hdr.Sync()
	if err != nil {
		return err
	}
	_, err = f.hdr.Seek(0, 0)
	if err != nil {
		return err
	}
	err = tarfile(f.hdr)
	if err != nil {
		return err
	}
	for _, t := range f.tuples {
		err = tarfile(t.f)
		if err != nil {
			return err
		}
	}
	err = tw.Close()
	if err != nil {
		return err
	}

	// clean-up
	err = f.hdr.Close()
	if err != nil {
		return err
	}
	err = os.RemoveAll(f.dir)
	if err != nil {
		return err
	}
	return err
}

// metadata about a HEP n-tuple
type metadata struct {
	Name     string
	Nentries int64
}

func make_metadata_from(t *Tuple) metadata {
	m := metadata{Name: t.f.Name(), Nentries: t.nentries}
	return m
}

// a HEP n-tuple
type Tuple struct {
	f        *os.File
	w        *writerz
	r        *readerz
	nentries int64
}

func (t *Tuple) Write(v interface{}) error {
	err := t.w.enc.Encode(v)
	if err == nil {
		t.nentries += 1
	}
	return err
}

func (t *Tuple) Close() error {
	//println("... closing ["+t.f.Name()+"] ...")
	var err error = nil
	if t.w != nil {
		err = t.w.compr.Close()
		if err != nil {
			return err
		}
	}
	if t.r != nil {
		err = t.r.compr.Close()
		if err != nil {
			return err
		}
	}
	err = t.f.Sync()
	if err != nil {
		return err
	}
	return t.f.Close()
}

func (t *Tuple) Entries() int64 {
	return t.nentries
}

func (t *Tuple) Read(v interface{}) error {
	err := t.r.dec.Decode(v)
	return err
}

// EOF
