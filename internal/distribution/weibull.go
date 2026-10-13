package distribution

import (
	"math"
	"math/rand/v2"
)

// weibull is the two-parameter Weibull distribution.
type weibull struct {
	shape float64 // β
	scale float64 // η
}

func (d weibull) Sample(rng *rand.Rand) float64 {
	return d.scale * math.Pow(drawPositive(rng.ExpFloat64), 1/d.shape)
}

func (d weibull) Mean() float64 { return d.scale * math.Gamma(1+1/d.shape) }

func (d weibull) CDF(t float64) float64 {
	if t <= 0 {
		return 0
	}
	return 1 - math.Exp(-math.Pow(t/d.scale, d.shape))
}

func (d weibull) PDF(t float64) float64 {
	if t < 0 {
		return 0
	}
	x := t / d.scale
	return d.shape / d.scale * math.Pow(x, d.shape-1) * math.Exp(-math.Pow(x, d.shape))
}
