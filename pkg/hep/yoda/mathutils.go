package yoda

import (
	"math"
)

// Compare 2 floating point numbers for equality with a degree of fuzziness
// The tolreance parameter is fractional.
func fFuzzyEq(a, b float64) bool {
	return fFuzzyEqWithTolerance(a, b, 1e-5)
}

// Compare 2 floating point numbers for equality with a degree of fuzziness
// The tolreance parameter is fractional.
func fFuzzyEqWithTolerance(a, b, tolerance float64) bool {
	absavg := math.Abs(a+b) / 2.0
	absdiff := math.Abs(a - b)
	return (absavg == 0.0 && absdiff == 0.0) || (absdiff < tolerance*absavg)
}

// Returns a list of nbins+1 values equally spaced between 'start' and 'end' inclusive
func linspace(start, end float64, nbins int) []float64 {
	if end < start {
		//FIXME(binet) proper error handling
		panic("invalid range")
	}
	if nbins <= 0 {
		//FIXME(binet) proper error handling
		panic("invalid range")
	}
	o := make([]float64, 0, nbins+1)
	interval := (end - start) / float64(nbins)
	for edge := start; edge <= end; edge += interval {
		o = append(o, edge)
	}
	if len(o) != nbins+1 {
		//FIXME(binet) proper error handling
		panic("internal error")
	}
	return o
}

/// Calculates the mean of a sample
func mean(sample []int) float64 {
	m := 0.0
	for _, i := range sample {
		m += float64(i)
	}
	return m / float64(len(sample))
}

/// Calculates the covariance (variance) between two samples
func covariance(s1, s2 []int) float64 {
	m1 := mean(s1)
	m2 := mean(s2)
	n := len(s1)
	if n > len(s2) {
		n = len(s2)
	}
	cov := 0.0
	for i := 0; i < n; i++ {
		c := (float64(s1[i]) - m1) * (float64(s2[i]) - m2)
		cov += c
	}
	if n > 1 {
		return cov / float64(n-1)
	}
	return 0.0
}

/// Calculates the correlation strength between two samples
func correlation(s1, s2 []int) float64 {
	//FIXME: make this cache-friendly by iterating only ONCE!
	cov := covariance(s1, s2)
	var1 := covariance(s1, s1)
	var2 := covariance(s2, s2)
	corr := cov / math.Sqrt(var1*var2)
	return corr * math.Sqrt(var2/var1)
}
