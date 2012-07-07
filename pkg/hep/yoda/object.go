package yoda

type Annotations map[string]interface{}

type Object interface {
	// Reset this analysis object
	Reset()

	// Retrieve the annotations attached to this analysis object
	Annotations() Annotations

	// Check if an annotation 'name' is defined
	HasAnnotation(name string) bool
}

type Bin interface {
	// Reset this bin
	Reset()

	// NumEntries returns the number of entries for this bin
	NumEntries() uint64

	// SumW returns the sum of weights
	SumW() float64

	// SumW2 returns the sum of weights squared
	SumW2() float64
}

type obj_impl struct {
	ann Annotations
}

