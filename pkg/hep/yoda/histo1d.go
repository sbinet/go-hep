package yoda

import (
	"errors"
	"math"
)

type hbin1d struct {
	Bin1D
}

// Fill this bin with weight 'weight' at position 'coord'
func (h *hbin1d) fill(coord, weight float64) {
	h.Bin1D.xdbn.fill(coord, weight)
}

// Fill this bin with weight 'weight'
func (h *hbin1d) fillBin(weight float64) {
	h.Bin1D.xdbn.fill(h.Bin1D.MidPoint(), weight)
}

// The area of a bin is the sum of weights in the bin, i.e. the
// width of the bin has no influence on this figure
func (h *hbin1d) area() float64 {
	return h.Bin1D.SumW()
}

// The height is defined as area/width
func (h *hbin1d) height() float64 {
	return h.area() / h.Bin1D.Width()
}

// Error computed using binomial statistics on the sum of bin weights
// i.e. err_area = sqrt{sum{weights}}
func (h *hbin1d) areaError() float64 {
	return math.Sqrt(h.Bin1D.SumW())
}

// Error includes scaling factor of the bin width
// i.e. err_height = sqrt{sum{weight}} / width
func (h *hbin1d) heightError() float64 {
	return h.areaError() / h.Bin1D.Width()
}

// Returns a new bin, the sum of the input bins
// if the edges do not match, nil is returned
func hbin1d_add(a, b *hbin1d) *hbin1d {
	if a.Bin1D.edges[0] != b.Bin1D.edges[0] ||
		a.Bin1D.edges[1] != b.Bin1D.edges[1] {
		return nil
	}
	return &hbin1d{Bin1D{edges: a.Bin1D.edges, xdbn: *dbn1d_add(&a.Bin1D.xdbn, &b.Bin1D.xdbn)}}
}

// in-place sum of the input bins: a += b
func hbin1d_iadd(a, b *hbin1d) error {
	if a.Bin1D.edges[0] != b.Bin1D.edges[0] ||
		a.Bin1D.edges[1] != b.Bin1D.edges[1] {
		return errors.New("iadd: axes' edges to not match")
	}
	return dbn1d_iadd(&a.Bin1D.xdbn, &b.Bin1D.xdbn)
}

// Returns a new bin, the subtraction of the input bins
func hbin1d_sub(a, b *hbin1d) *hbin1d {
	if a.Bin1D.edges[0] != b.Bin1D.edges[0] ||
		a.Bin1D.edges[1] != b.Bin1D.edges[1] {
		return nil
	}
	return &hbin1d{Bin1D{edges: a.Bin1D.edges, xdbn: *dbn1d_sub(&a.Bin1D.xdbn, &b.Bin1D.xdbn)}}
}

// Compares 2 hbin1ds, by lower edge position
// FIXME: check for overlap somewhere...
func hbin1d_less(a, b *hbin1d) bool {
	return a.Bin1D.edges[0] < b.Bin1D.edges[0]
}

// a list of sorted hbin1d
type sorted_hbin1ds []hbin1d

func (s sorted_hbin1ds) Len() int {
	return len(s)
}

func (s sorted_hbin1ds) Less(i, j int) bool {
	return hbin1d_less(&s[i], &s[j])
}

func (s sorted_hbin1ds) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
// A one-dimensional histogram
type Histo1D struct {
	axis Axis1D
}
