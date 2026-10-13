// Package analysis holds the application's use cases: it turns each kind of
// input data into reliability indicators ready to be presented.
package analysis

import (
	"math/rand/v2"
	"time"

	"empirical-reliability/internal/distribution"
	"empirical-reliability/internal/reliability"
)

// seedStreamMix derives the second word of the generator state from the seed,
// so that a single number reproduces a sample.
const seedStreamMix = 0x9e3779b97f4a7c15

// Analysis is the outcome of one calculation.
type Analysis struct {
	Result *reliability.Result
	// Random describes the drawn sample; nil for manually entered data.
	Random *RandomSample
}

// RandomSample records how a random sample was drawn.
type RandomSample struct {
	Spec       distribution.Spec
	Parameters []float64
	Seed       uint64
	Model      distribution.Distribution
}

// RandomRequest describes a random sample to draw and analyse.
type RandomRequest struct {
	Spec       distribution.Spec
	Parameters []float64
	Size       int
	Seed       uint64
	Step       reliability.StepSelector
}

// Model returns the theoretical distribution the sample was drawn from, or nil
// for manually entered data.
func (a *Analysis) Model() distribution.Model {
	if a.Random == nil {
		return nil
	}
	return a.Random.Model
}

// FromTable analyses grouped data: a sample size and a failure table.
func FromTable(sampleSize int, observations []reliability.Observation) (*Analysis, error) {
	result, err := reliability.Compute(sampleSize, observations)
	if err != nil {
		return nil, err
	}
	return &Analysis{Result: result}, nil
}

// FromTimes analyses individual failure times.
func FromTimes(times []float64, sampleSize int, step reliability.StepSelector) (*Analysis, error) {
	result, err := reliability.FromTimes(times, sampleSize, step)
	if err != nil {
		return nil, err
	}
	return &Analysis{Result: result}, nil
}

// FromRandom draws a sample from the requested distribution and analyses it.
// The same request always yields the same sample.
func FromRandom(request RandomRequest) (*Analysis, error) {
	model, err := request.Spec.New(request.Parameters...)
	if err != nil {
		return nil, err
	}
	rng := rand.New(rand.NewPCG(request.Seed, request.Seed^seedStreamMix))
	times, err := distribution.Sample(model, request.Size, rng)
	if err != nil {
		return nil, err
	}
	result, err := reliability.FromTimes(times, len(times), request.Step)
	if err != nil {
		return nil, err
	}
	return &Analysis{
		Result: result,
		Random: &RandomSample{
			Spec:       request.Spec,
			Parameters: request.Parameters,
			Seed:       request.Seed,
			Model:      model,
		},
	}, nil
}

// TimeSeed returns a seed that differs between runs.
func TimeSeed() uint64 {
	return uint64(time.Now().UnixNano())
}
