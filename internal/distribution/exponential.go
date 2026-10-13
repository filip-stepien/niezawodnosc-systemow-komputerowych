package distribution

import (
	"math"
	"math/rand/v2"
)

// exponential is the exponential distribution with a constant hazard rate.
type exponential struct {
	rate float64 // λ
}

func (d exponential) Sample(rng *rand.Rand) float64 {
	return drawPositive(rng.ExpFloat64) / d.rate
}

func (d exponential) Mean() float64 { return 1 / d.rate }

func (d exponential) CDF(t float64) float64 {
	if t <= 0 {
		return 0
	}
	return 1 - math.Exp(-d.rate*t)
}

func (d exponential) PDF(t float64) float64 {
	if t < 0 {
		return 0
	}
	return d.rate * math.Exp(-d.rate*t)
}
