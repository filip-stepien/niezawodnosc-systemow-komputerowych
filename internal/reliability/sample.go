package reliability

import "slices"

// FromTimes estimates the indicators from individual failure times, grouping
// them into intervals whose length is chosen by selectStep. When sampleSize
// exceeds the number of times, the remaining elements are treated as still
// working at the end of the observation.
func FromTimes(times []float64, sampleSize int, selectStep StepSelector) (*Result, error) {
	if err := validateTimes(times); err != nil {
		return nil, err
	}
	if sampleSize < len(times) {
		return nil, ErrSampleSmallerThanTimes
	}
	observations, err := groupValidated(times, selectStep(times))
	if err != nil {
		return nil, err
	}
	result, err := Compute(sampleSize, observations)
	if err != nil {
		return nil, err
	}

	result.Times = slices.Sorted(slices.Values(times))
	if result.Survived == 0 {
		result.Mean = mean(times)
	}
	return result, nil
}

func mean(values []float64) float64 {
	total := 0.0
	for _, v := range values {
		total += v
	}
	return total / float64(len(values))
}
