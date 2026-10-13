// Package distribution provides the theoretical lifetime distributions a
// random sample can be drawn from and compared against.
package distribution

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
)

// MaxSampleSize is the largest random sample that can be drawn.
const MaxSampleSize = 1_000_000

// negligibleSurvival is the reliability below which the hazard rate is no
// longer computed, because the quotient becomes numerically meaningless.
const negligibleSurvival = 1e-12

var (
	ErrSampleSizeNotPositive = errors.New("sample size must be positive")
	ErrSampleSizeTooLarge    = errors.New("sample size is too large")
	ErrParameterCount        = errors.New("wrong number of distribution parameters")
)

// ParameterError reports a distribution parameter that is not positive.
type ParameterError struct {
	Symbol string
}

func (e *ParameterError) Error() string {
	return fmt.Sprintf("parameter %s must be positive", e.Symbol)
}

// Sampler draws lifetimes.
type Sampler interface {
	Sample(rng *rand.Rand) float64
}

// Model describes the theoretical shape of a lifetime distribution.
type Model interface {
	CDF(t float64) float64 // F(t)
	PDF(t float64) float64 // f(t)
}

// Distribution is a lifetime distribution of the random variable T.
type Distribution interface {
	Sampler
	Model
	Mean() float64 // ET
}

// Hazard returns the theoretical hazard rate λ(t) = f(t) / (1 - F(t)), or NaN
// where the reliability is negligible.
func Hazard(m Model, t float64) float64 {
	survival := 1 - m.CDF(t)
	if survival < negligibleSurvival {
		return math.NaN()
	}
	return m.PDF(t) / survival
}

// Sample draws size lifetimes from s.
func Sample(s Sampler, size int, rng *rand.Rand) ([]float64, error) {
	if size <= 0 {
		return nil, ErrSampleSizeNotPositive
	}
	if size > MaxSampleSize {
		return nil, ErrSampleSizeTooLarge
	}
	times := make([]float64, size)
	for i := range times {
		times[i] = s.Sample(rng)
	}
	return times, nil
}

// drawPositive repeats draw until it yields a positive value, because a
// lifetime of zero is not admissible.
func drawPositive(draw func() float64) float64 {
	for {
		if value := draw(); value > 0 {
			return value
		}
	}
}

func isPositive(value float64) bool {
	return value > 0 && !math.IsInf(value, 0)
}
