package reliability

import (
	"math"
	"slices"
)

// MaxIntervals limits the number of intervals failure times may be grouped into.
const MaxIntervals = 10000

const (
	minAutoIntervals = 5
	maxAutoIntervals = 50

	// boundaryTolerance keeps a time lying on an interval boundary in the
	// lower interval despite floating-point rounding.
	boundaryTolerance = 1e-9
)

// roundMultipliers are the leading digits a rounded interval length may have.
var roundMultipliers = []float64{1, 2, 2.5, 5, 10}

// StepSelector chooses the interval length Δt for grouping failure times.
type StepSelector func(times []float64) float64

// FixedStep returns a selector that always chooses the given interval length.
func FixedStep(step float64) StepSelector {
	return func([]float64) float64 { return step }
}

// AutoStep chooses the interval length from the data: about √n intervals,
// with the length rounded up to a round number (1, 2, 2.5 or 5 times a power
// of ten). It is a StepSelector.
func AutoStep(times []float64) float64 {
	if len(times) == 0 {
		return math.NaN()
	}
	intervals := math.Ceil(math.Sqrt(float64(len(times))))
	intervals = min(max(intervals, minAutoIntervals), maxAutoIntervals)
	return roundUp(slices.Max(times) / intervals)
}

func roundUp(value float64) float64 {
	magnitude := math.Pow(10, math.Floor(math.Log10(value)))
	for _, multiplier := range roundMultipliers {
		if rounded := multiplier * magnitude; rounded >= value*(1-boundaryTolerance) {
			return rounded
		}
	}
	return value
}

// Group counts failure times in the intervals (t_{i-1}, t_i], where t_i = i·step.
func Group(times []float64, step float64) ([]Observation, error) {
	if err := validateTimes(times); err != nil {
		return nil, err
	}
	return groupValidated(times, step)
}

func groupValidated(times []float64, step float64) ([]Observation, error) {
	if !isFinite(step) || step <= 0 {
		return nil, ErrStepNotPositive
	}
	longest := slices.Max(times)
	if longest/step > MaxIntervals {
		return nil, ErrTooManyIntervals
	}

	observations := make([]Observation, intervalIndex(longest, step))
	for i := range observations {
		observations[i].Time = float64(i+1) * step
	}
	for _, t := range times {
		observations[intervalIndex(t, step)-1].Failures++
	}
	return observations, nil
}

func validateTimes(times []float64) error {
	if len(times) == 0 {
		return ErrNoFailureTimes
	}
	for i, t := range times {
		if !isFinite(t) || t <= 0 {
			return &ItemError{Position: i + 1, Err: ErrFailureTimeNotPositive}
		}
	}
	return nil
}

// intervalIndex returns the 1-based index i of the interval (t_{i-1}, t_i]
// that contains t.
func intervalIndex(t, step float64) int {
	return max(int(math.Ceil(t/step-boundaryTolerance)), 1)
}
