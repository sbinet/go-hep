package fourmom

type I4Mom interface {
	Px() float64
}

type PxPyPzE struct {
	mom [4]float64
}

// EOF
