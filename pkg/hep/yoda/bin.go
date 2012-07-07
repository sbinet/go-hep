package yoda

import (
	"errors"
	"math"

	//"sort"
)

// Base class for bins in 1D normal and profile histograms.
// The lower bin edge is inclusive
type Bin1D struct {
	// the bin limits
	edges [2]float64

	// distribution of weighted x-values
	xdbn dbn1d
}

// Create a new Bin1D given low and high edges
func NewBin1D(low, hi float64) *Bin1D {
	return &Bin1D{edges: [2]float64{low, hi}, xdbn: dbn1d{}}
}

// XMin returns the lower limit of the bin (inclusive)
func (b *Bin1D) XMin() float64 {
	return b.edges[0]
}

// Returns the lower limit of the bin (inclusive)
func (b *Bin1D) LowEdge() float64 {
	return b.edges[0]
}

// XMax returns the upper limit of the bin (exclusive)
func (b *Bin1D) XMax() float64 {
	return b.edges[1]
}

// Returns the upper limit of the bin (exclusive)
func (b *Bin1D) HighEdge() float64 {
	return b.edges[1]
}

// Edges returns the lower and upper limits of the bin
func (b *Bin1D) Edges() (float64, float64) {
	return b.edges[0], b.edges[1]
}

// Reset this object
func (b *Bin1D) Reset() {
	b.xdbn.Reset()
}

// Returns the separation of low and high edges, i,e, high - low
func (b *Bin1D) Width() float64 {
	return b.edges[1] - b.edges[0]
}

// Returns the mean position in the bin, or the midpoint if that is not available
func (b *Bin1D) Focus() float64 {
	if b.xdbn.sumw != 0.0 {
		return b.XMean()
	}
	return b.MidPoint()
}

// Returns the geometric centre of the bin, i.e. high+low/2.0
func (b *Bin1D) MidPoint() float64 {
	return (b.edges[1] + b.edges[0]) / 2.0
}

// Returns the mean value of x-values in the bin.
func (b *Bin1D) XMean() float64 {
	return b.xdbn.mean()
}

// Returns the variance of x-values in the bin
func (b *Bin1D) XVariance() float64 {
	v, err := b.xdbn.variance()
	if err != nil {
		panic(err)
	}
	return v
}

// Returns the standard deviation (spread) of x-values in the bin
func (b *Bin1D) XStdDev() float64 {
	v, err := b.xdbn.stdDev()
	if err != nil {
		panic(err)
	}
	return v
}

// Returns the standard error on the bin focus
func (b *Bin1D) XStdError() float64 {
	v, err := b.xdbn.stdErr()
	if err != nil {
		panic(err)
	}
	return v
}

// Returns the number of entries in the bin
func (b *Bin1D) NumEntries() uint64 {
	return b.xdbn.nfills
}

// Returns the sum of weights
func (b *Bin1D) SumW() float64 {
	return b.xdbn.sumw
}

// Returns the sum of weights squared
func (b *Bin1D) SumW2() float64 {
	return b.xdbn.sumw2
}

// Returns the sum of x*weight
func (b *Bin1D) SumWX() float64 {
	return b.xdbn.sumwx
}

// Returns the sum of x**2 * weight
func (b *Bin1D) SumWX2() float64 {
	return b.xdbn.sumwx2
}

func (b *Bin1D) scaleW(scale float64) {
	b.xdbn.scaleW(scale)
}

// Returns a new bin, the sum of the input bins
// if the edges do not match, nil is returned
func Bin1D_Add(a, b *Bin1D) *Bin1D {
	if a.edges[0] != b.edges[0] ||
		a.edges[1] != b.edges[1] {
		return nil
	}
	return &Bin1D{edges: a.edges, xdbn: *dbn1d_add(&a.xdbn, &b.xdbn)}
}

// in-place sum of the input bins: a += b
func Bin1D_IAdd(a, b *Bin1D) error {
	if a.edges[0] != b.edges[0] ||
		a.edges[1] != b.edges[1] {
		return errors.New("iadd: axes' edges to not match")
	}
	return dbn1d_iadd(&a.xdbn, &b.xdbn)
}

// Returns a new bin, the subtraction of the input bins
func Bin1D_Sub(a, b *Bin1D) *Bin1D {
	if a.edges[0] != b.edges[0] ||
		a.edges[1] != b.edges[1] {
		return nil
	}
	return &Bin1D{edges: a.edges, xdbn: *dbn1d_sub(&a.xdbn, &b.xdbn)}
}

// Compares 2 Bin1Ds, by lower edge position
// FIXME: check for overlap somewhere...
func Bin1D_Less(a, b *Bin1D) bool {
	return a.edges[0] < b.edges[0]
}

// dbn1d is a 1D distribution.
// This class is used to centralize the calculation of statistics of unbounded,
// unbinned sampled distributions.
// Each distribution fill contributes a weight 'w' and a value 'x'.
// By storing the total number of fills (ignoring weights), Sum(w), Sum(w^2),
// Sum(w.x) and Sum(w.x^2), the dbn1d can calculate the mean and spread 
// \sigma^2, \sigma and \hat{\sigma} of the sampled distribution.
// It is used to provide this information in bins for the "hidden" 'y' 
// distribution in profile histogram bins.
type dbn1d struct {
	nfills uint64
	sumw   float64
	sumw2  float64
	sumwx  float64
	sumwx2 float64
}

func (d *dbn1d) scaleW(factor float64) {
	sf2 := factor * factor
	d.sumw *= factor
	d.sumw2 *= sf2
	d.sumwx *= factor
	d.sumwx2 *= sf2
}

func (d *dbn1d) fill(val, weight float64) {
	d.nfills += 1
	d.sumw += weight
	w2 := weight * weight
	if weight < 0. {
		w2 *= -1.0
	}
	d.sumw2 += w2
	d.sumwx += weight * val
	d.sumwx2 += weight * val * val
}

func (d *dbn1d) Reset() {
	d.nfills = 0
	d.sumw = 0.0
	d.sumw2 = 0.0
	d.sumwx = 0.0
	d.sumwx2 = 0.0
}

func (d *dbn1d) effNumEntries() float64 {
	return d.sumw * d.sumw / d.sumw2
}

func (d *dbn1d) mean() float64 {
	return d.sumwx / d.sumw
}

// The weighted variance is defined as:
//  sig2 = (sum(wx**2) * sum(w) - sum(wx)**2) / (sum(w)**2 - sum(w**2))
//  http://en.wikipedia.org/wiki/Weighted_mean
func (d *dbn1d) variance() (float64, error) {
	effn := d.effNumEntries()
	if effn == 0.0 {
		return 0.0, errors.New("requested width of a distribution with no net fill weights")
	}
	if effn <= 1.0 {
		return 0.0, errors.New("requested width of a distribution with only one effective entry")
	}
	num := d.sumwx2*d.sumw - d.sumwx*d.sumwx
	den := d.sumw*d.sumw - d.sumw2
	if den == 0 {
		panic("undefined weighted variance")
	}
	if math.Abs(num) < 1e-10 && math.Abs(den) < 1e-10 {
		panic("numerically unstable weights in width calculation")
	}
	return num / den, nil
}

func (d *dbn1d) stdDev() (float64, error) {
	v, err := d.variance()
	if err != nil {
		return 0.0, err
	}
	return math.Sqrt(v), nil
}

func (d *dbn1d) stdErr() (float64, error) {
	effnum := d.effNumEntries()
	if effnum == 0.0 {
		return 0.0, errors.New("requested std error of a distribution with no net fill weights")
	}
	v, err := d.variance()
	if err != nil {
		return 0.0, err
	}
	return math.Sqrt(v / effnum), nil
}

func dbn1d_add(a, b *dbn1d) *dbn1d {
	return &dbn1d{
		nfills: a.nfills + b.nfills,
		sumw:   a.sumw + b.sumw,
		sumw2:  a.sumw2 + b.sumw2,
		sumwx:  a.sumwx + b.sumwx,
		sumwx2: a.sumwx2 + b.sumwx2,
	}
}

func dbn1d_iadd(a, b *dbn1d) error {
	a.nfills += b.nfills
	a.sumw += b.sumw
	a.sumw2 += b.sumw2
	a.sumwx += b.sumwx
	a.sumwx2 += b.sumwx2
	return nil
}

func dbn1d_sub(a, b *dbn1d) *dbn1d {
	return &dbn1d{
		nfills: a.nfills + b.nfills,
		sumw:   a.sumw - b.sumw,
		sumw2:  a.sumw2 - b.sumw2,
		sumwx:  a.sumwx - b.sumwx,
		sumwx2: a.sumwx2 - b.sumwx2,
	}
}

func dbn1d_isub(a, b *dbn1d) error {
	a.nfills += b.nfills
	a.sumw -= b.sumw
	a.sumw2 -= b.sumw2
	a.sumwx -= b.sumwx
	a.sumwx2 -= b.sumwx2
	return nil
}

// a list of sorted Bin1Ds
type SortedBin1Ds []Bin1D

func (s SortedBin1Ds) Len() int {
	return len(s)
}

func (s SortedBin1Ds) Less(i, j int) bool {
	return Bin1D_Less(&s[i], &s[j])
}

func (s SortedBin1Ds) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
