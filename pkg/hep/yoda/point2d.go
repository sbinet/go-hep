package yoda

// A 2D data point to be contained in a Scatter2D
type Point2D struct {
	coord [2]float64
	err   [2][2]float64
}

func NewPoint2D(x, y float64) *Point2D {
	return NewPoint2DErr(x, y, 0., 0.)
}

func NewPoint2DErr(x, y, ex, ey float64) *Point2D {
	return &Point2D{
		coord: [2]float64{x, y},
		err:   [2][2]float64{{ex, ex}, {ey, ey}},
	}
}

func NewPoint2DAsymErr(x, y, exmin, exmax, eymin, eymax float64) *Point2D {
	return &Point2D{
		coord: [2]float64{x, y},
		err:   [2][2]float64{{exmin, exmax}, {eymin, eymax}},
	}
}

// Get x-coordinate
func (p *Point2D) X() float64 {
	return p.coord[0]
}

// Get x-value minus negative x-error
func (p *Point2D) XMin() float64 {
	return p.coord[0] - p.err[0][0]
}

// Get x-value plus positive x-error
func (p *Point2D) XMax() float64 {
	return p.coord[0] + p.err[0][1]
}

// Get x-error values
func (p *Point2D) XErrs() (float64, float64) {
	return p.err[0][0], p.err[0][1]
}

// Get negative x-error value
func (p *Point2D) XErrMinus() float64 {
	return p.err[0][0]
}

// Get positive x-error value
func (p *Point2D) XErrPlus() float64 {
	return p.err[0][1]
}

// Get average x-error value
func (p *Point2D) XErrAvg() float64 {
	return (p.err[0][0] + p.err[0][1]) / float64(2.)
}

// Set x-coordinate
func (p *Point2D) SetX(x float64) {
	p.coord[0] = x
}

// Set symmetric x-error
func (p *Point2D) SetXErr(ex float64) {
	p.err[0][0] = ex
	p.err[0][1] = ex
}

// Set asymmetric x-error
func (p *Point2D) SetXErrs(exmin, exmax float64) {
	p.err[0][0] = exmin
	p.err[0][1] = exmax
}

// Get y-coordinate
func (p *Point2D) Y() float64 {
	return p.coord[1]
}

// Get y-value minus negative y-error
func (p *Point2D) YMin() float64 {
	return p.coord[1] - p.err[1][0]
}

// Get y-value plus positive y-error
func (p *Point2D) YMax() float64 {
	return p.coord[1] + p.err[1][1]
}

// Get y-error values
func (p *Point2D) YErrs() (float64, float64) {
	return p.err[1][0], p.err[1][1]
}

// Get negative y-error value
func (p *Point2D) YErrMinus() float64 {
	return p.err[1][0]
}

// Get positive y-error value
func (p *Point2D) YErrPlus() float64 {
	return p.err[1][1]
}

// Get average y-error value
func (p *Point2D) YErrAvg() float64 {
	return (p.err[1][0] + p.err[1][1]) / float64(2.)
}

// Set y-coordinate
func (p *Point2D) SetY(y float64) {
	p.coord[1] = y
}

// Set symmetric y-error
func (p *Point2D) SetYErr(ey float64) {
	p.err[1][0] = ey
	p.err[1][1] = ey
}

// Set asymmetric y-error
func (p *Point2D) SetYErrs(eymin, eymax float64) {
	p.err[1][0] = eymin
	p.err[1][1] = eymax
}

/// Equality test (of x-characteristics only)
func Point2D_Eq(p1, p2 *Point2D) bool {
	return fFuzzyEq(p1.X(), p2.X()) &&
		fFuzzyEq(p1.XErrMinus(), p2.XErrMinus()) &&
		fFuzzyEq(p1.XErrPlus(), p2.XErrPlus())
}
