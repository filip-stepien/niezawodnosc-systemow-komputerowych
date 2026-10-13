package reliability

import (
	"errors"
	"fmt"
)

// Validation errors reported for malformed input data.
var (
	ErrSampleSizeNotPositive  = errors.New("sample size must be positive")
	ErrNoObservations         = errors.New("no observations")
	ErrTimeNotIncreasing      = errors.New("observation times must be positive and strictly increasing")
	ErrNegativeFailures       = errors.New("failure count must not be negative")
	ErrFailuresExceedSample   = errors.New("failures exceed the sample size")
	ErrNoFailureTimes         = errors.New("no failure times")
	ErrFailureTimeNotPositive = errors.New("failure time must be positive")
	ErrStepNotPositive        = errors.New("interval length must be positive")
	ErrTooManyIntervals       = errors.New("interval length yields too many intervals")
	ErrSampleSmallerThanTimes = errors.New("sample size is smaller than the number of failure times")
)

// ItemError locates an error at a 1-based position in a sequence of input values.
type ItemError struct {
	Position int
	Err      error
}

func (e *ItemError) Error() string { return fmt.Sprintf("item %d: %v", e.Position, e.Err) }
func (e *ItemError) Unwrap() error { return e.Err }
