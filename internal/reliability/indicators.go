// Package reliability computes the empirical reliability indicators F*(t),
// R*(t), f*(t), λ*(t) and E*T as defined in: Niezawodność urządzeń
// elektronicznych, ed. M. Prażewska, WKiŁ, Warszawa 1987, pp. 208-213.
// Equation numbers in comments refer to that book.
package reliability

import "math"

// initialReliability is R(0): every element works at t = 0.
const initialReliability = 1.0

// Observation is one inspection: Failures elements were found failed at Time,
// i.e. they failed in the interval (previous inspection, Time].
type Observation struct {
	Time     float64
	Failures int
}

// Row holds the indicators estimated for one inspection time t_i.
type Row struct {
	Time        float64 // t_i
	Step        float64 // Δt = t_i - t_{i-1}
	Failures    int     // m(t_i)
	Working     int     // n(t_i), elements still working at t_i
	CDF         float64 // F*(t_i), eq. (5-2)
	Reliability float64 // R*(t_i), eq. (5-5)
	Density     float64 // f*(t_i), eq. (5-9)
	Hazard      float64 // λ*(t_i), eq. (5-12); NaN when no element works at the start of the interval
}

// Result is the complete set of indicators for one sample.
type Result struct {
	SampleSize int // n(0)
	Failed     int
	Survived   int // elements still working when the observation ended
	Rows       []Row

	// Mean is E*T from eq. (5-15). Grouped data carry no individual lifetimes,
	// so the inspection time is used as the lifetime of every element found
	// failed at it. NaN when some elements survived: their lifetimes are unknown.
	Mean float64
	// MeanFromReliability is E*T from eq. (5-13), integrating R*(t) with the
	// trapezoidal rule. It is a lower bound when some elements survived.
	MeanFromReliability float64

	// Times holds the individual failure times in ascending order; empty for
	// grouped data.
	Times []float64
}

// Compute estimates the indicators from grouped data: the sample size n(0) and
// the failure counts found at consecutive inspection times.
func Compute(sampleSize int, observations []Observation) (*Result, error) {
	if err := validateObservations(sampleSize, observations); err != nil {
		return nil, err
	}
	rows := buildRows(sampleSize, observations)
	survived := rows[len(rows)-1].Working
	return &Result{
		SampleSize:          sampleSize,
		Failed:              sampleSize - survived,
		Survived:            survived,
		Rows:                rows,
		Mean:                meanOfInspectionTimes(sampleSize, survived, observations),
		MeanFromReliability: integrateReliability(rows),
	}, nil
}

func validateObservations(sampleSize int, observations []Observation) error {
	if sampleSize <= 0 {
		return ErrSampleSizeNotPositive
	}
	if len(observations) == 0 {
		return ErrNoObservations
	}
	previousTime, working := 0.0, sampleSize
	for i, o := range observations {
		if err := validateObservation(o, previousTime, working); err != nil {
			return &ItemError{Position: i + 1, Err: err}
		}
		previousTime, working = o.Time, working-o.Failures
	}
	return nil
}

func validateObservation(o Observation, previousTime float64, working int) error {
	switch {
	case !isFinite(o.Time) || o.Time <= previousTime:
		return ErrTimeNotIncreasing
	case o.Failures < 0:
		return ErrNegativeFailures
	case o.Failures > working:
		return ErrFailuresExceedSample
	}
	return nil
}

func buildRows(sampleSize int, observations []Observation) []Row {
	rows := make([]Row, 0, len(observations))
	previousTime, working := 0.0, sampleSize
	for _, o := range observations {
		row := newRow(o, o.Time-previousTime, working, sampleSize)
		rows = append(rows, row)
		previousTime, working = o.Time, row.Working
	}
	return rows
}

// newRow estimates the indicators for one interval of length step that starts
// with workingAtStart working elements.
func newRow(o Observation, step float64, workingAtStart, sampleSize int) Row {
	working := workingAtStart - o.Failures
	return Row{
		Time:        o.Time,
		Step:        step,
		Failures:    o.Failures,
		Working:     working,
		CDF:         float64(sampleSize-working) / float64(sampleSize),
		Reliability: float64(working) / float64(sampleSize),
		Density:     float64(o.Failures) / (float64(sampleSize) * step),
		Hazard:      hazard(o.Failures, workingAtStart, step),
	}
}

func hazard(failures, workingAtStart int, step float64) float64 {
	if workingAtStart == 0 {
		return math.NaN()
	}
	return float64(failures) / (float64(workingAtStart) * step)
}

func meanOfInspectionTimes(sampleSize, survived int, observations []Observation) float64 {
	if survived > 0 {
		return math.NaN()
	}
	total := 0.0
	for _, o := range observations {
		total += o.Time * float64(o.Failures)
	}
	return total / float64(sampleSize)
}

func integrateReliability(rows []Row) float64 {
	area, previous := 0.0, initialReliability
	for _, row := range rows {
		area += (previous + row.Reliability) / 2 * row.Step
		previous = row.Reliability
	}
	return area
}

func isFinite(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0)
}
