package distribution

import (
	"math"
	"math/rand/v2"
)

// normal is the normal distribution truncated to positive values, because a
// lifetime cannot be negative.
type normal struct {
	mean      float64 // μ of the untruncated distribution
	deviation float64 // σ of the untruncated distribution
	mass      float64 // probability mass left after the truncation
}

func newNormal(mean, deviation float64) normal {
	return normal{
		mean:      mean,
		deviation: deviation,
		mass:      1 - standardCDF(-mean/deviation),
	}
}

func (d normal) Sample(rng *rand.Rand) float64 {
	return drawPositive(func() float64 { return d.mean + d.deviation*rng.NormFloat64() })
}

func (d normal) Mean() float64 {
	return d.mean + d.deviation*standardPDF(-d.mean/d.deviation)/d.mass
}

func (d normal) CDF(t float64) float64 {
	if t <= 0 {
		return 0
	}
	return (standardCDF((t-d.mean)/d.deviation) - (1 - d.mass)) / d.mass
}

func (d normal) PDF(t float64) float64 {
	if t < 0 {
		return 0
	}
	return standardPDF((t-d.mean)/d.deviation) / (d.deviation * d.mass)
}

func standardCDF(x float64) float64 { return 0.5 * math.Erfc(-x/math.Sqrt2) }
func standardPDF(x float64) float64 { return math.Exp(-x*x/2) / math.Sqrt(2*math.Pi) }
