package yoda

import "sort"

// Axis1D is a container of ordered bins
type Axis1D struct {
	// the bins contained in this histogram
	bins []hbin1d
	// a distribution counter for underflow fills
	underflow dbn1d
	// a distribution counter for overflow fills
	overflow dbn1d

	// a distribution counter for the whole histogram
	dbn dbn1d

	// bin edges: lower edges, except last entry, which is the high edge of the last entry
	edges []float64

	// hash table for fast bin lookup
	hash map[float64]int
}

// Create a new Axis1D from a list of bin edges
func NewAxis1DFromEdges(edges []float64) *Axis1D {
	nbins := len(edges) - 1
	a := &Axis1D{bins: make([]hbin1d, 0, nbins),
		underflow: dbn1d{},
		overflow:  dbn1d{},
		dbn:       dbn1d{},
		edges:     make([]float64, 0, len(edges)),
		hash:      make(map[float64]int),
	}
	copy(a.edges, edges)
	for i := 0; i < nbins; i++ {
		a.bins = append(a.bins, hbin1d{*NewBin1D(edges[i], edges[i+1])})
	}
	sort.Sort((sorted_hbin1ds)(a.bins))
	for i := range a.bins {
		// insert upper bound mapped to bin id
		a.hash[edges[i+1]] = i
	}
	return a
}

// Create a new Axis1D from a number of bins and a bin distribution
func NewAxis1D(nbins int, lower, upper float64) *Axis1D {
	return NewAxis1DFromEdges(linspace(lower, upper, nbins))
}

// Returns the number of bins (not counting under|over-flows)
func (a *Axis1D) NumBins() uint64 {
	return uint64(len(a.bins))
}

// Returns the bins' axis
func (a *Axis1D) Bins() []hbin1d {
	return a.bins
}

// Returns the edges of the bin number 'id'
func (a *Axis1D) BinEdges(id int) (float64, float64) {
	return a.edges[id], a.edges[id+1]
}

// Returns the low edge of the axis
func (a *Axis1D) LowEdge() float64 {
	return a.bins[0].LowEdge()
}

// Returns the high edge of the axis
func (a *Axis1D) HighEdge() float64 {
	return a.bins[len(a.bins)-1].LowEdge()
}

// Returns the bin at bin number 'id'
func (a *Axis1D) Bin(id int) *hbin1d {
	return &a.bins[id]
}

// Returns the bin at coordinate 'x'
func (a *Axis1D) BinByCoord(x float64) *hbin1d {
	id, ok := a.hash[x]
	if ok {
		return &a.bins[id]
	}
	//FIXME: O(N) search
	for i := range a.edges[1:] {
		if a.edges[i-1] <= x && a.edges[i] < x {
			id = i - 1
			a.hash[x] = id
			return &a.bins[id]
		}
	}
	return nil
}

// Reset the axis content
func (a *Axis1D) Reset() {
	a.dbn.Reset()
	a.underflow.Reset()
	a.overflow.Reset()
	for i := range a.bins {
		a.bins[i].Reset()
	}
}

// Scale the axis coordinates
func (a *Axis1D) ScaleW(scale float64) {
	a.dbn.scaleW(scale)
	a.underflow.scaleW(scale)
	a.overflow.scaleW(scale)
	for i := range a.bins {
		a.bins[i].scaleW(scale)
	}
}

// In-place add of 2 axes
func Axis1D_IAdd(a, b *Axis1D) error {
	//FIXME error handling
	if len(a.bins) != len(b.bins) {
		panic("axes lengthes differ")
	}
	for i := range a.bins {
		err := hbin1d_iadd(&a.bins[i], &b.bins[i])
		if err != nil {
			return err
		}
	}
	return dbn1d_iadd(&a.dbn, &b.dbn)
}
