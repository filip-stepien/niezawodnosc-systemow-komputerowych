package reliability

import (
	"errors"
	"math"
	"testing"
)

const (
	bookSampleSize = 35
	// bookTolerance allows for the rounding of the tables printed in the book
	// (e.g. 0.828 instead of 0.829).
	bookTolerance = 0.0015
	exact         = 1e-12
)

// bookObservations is the data of example 5.1 (p. 209): 35 capacitors
// inspected every hour.
var bookObservations = []Observation{
	{1, 0}, {2, 3}, {3, 3}, {4, 5}, {5, 8}, {6, 7}, {7, 6}, {8, 2}, {9, 1}, {10, 0},
}

func near(got, want, tolerance float64) bool { return math.Abs(got-want) <= tolerance }

func bookResult(t *testing.T) *Result {
	t.Helper()
	result, err := Compute(bookSampleSize, bookObservations)
	if err != nil {
		t.Fatalf("Compute: %v", err)
	}
	return result
}

func TestComputeMatchesBookExample51CDF(t *testing.T) {
	want := []float64{0, 0.086, 0.171, 0.314, 0.543, 0.743, 0.914, 0.971, 1, 1}
	for i, row := range bookResult(t).Rows {
		if !near(row.CDF, want[i], bookTolerance) {
			t.Errorf("F*(t_%d) = %.4f, book gives %.3f", i+1, row.CDF, want[i])
		}
	}
}

func TestComputeMatchesBookExample52Reliability(t *testing.T) {
	want := []float64{1, 0.914, 0.828, 0.686, 0.457, 0.257, 0.086, 0.029, 0, 0}
	for i, row := range bookResult(t).Rows {
		if !near(row.Reliability, want[i], bookTolerance) {
			t.Errorf("R*(t_%d) = %.4f, book gives %.3f", i+1, row.Reliability, want[i])
		}
		if !near(row.CDF+row.Reliability, 1, exact) {
			t.Errorf("F*(t_%d) + R*(t_%d) = %v, want 1", i+1, i+1, row.CDF+row.Reliability)
		}
	}
}

func TestComputeMatchesBookExample53Density(t *testing.T) {
	want := []float64{0, 0.086, 0.086, 0.143, 0.228, 0.2, 0.171, 0.057, 0.0286, 0}
	area := 0.0
	for i, row := range bookResult(t).Rows {
		if !near(row.Density, want[i], bookTolerance) {
			t.Errorf("f*(t_%d) = %.4f, book gives %.4f", i+1, row.Density, want[i])
		}
		area += row.Density * row.Step
	}
	if !near(area, 1, exact) {
		t.Errorf("sum of f*(t_i)·Δt = %v, want 1", area)
	}
}

func TestComputeMatchesBookExample54Hazard(t *testing.T) {
	want := []float64{0, 0.086, 0.095, 0.172, 0.333, 0.437, 0.667, 0.667, 1.0}
	rows := bookResult(t).Rows
	for i, w := range want {
		if !near(rows[i].Hazard, w, bookTolerance) {
			t.Errorf("λ*(t_%d) = %.4f, book gives %.3f", i+1, rows[i].Hazard, w)
		}
	}
	if last := rows[len(rows)-1].Hazard; !math.IsNaN(last) {
		t.Errorf("λ*(t_10) = %v, want undefined: no element works at the start of the interval", last)
	}
}

func TestComputeMeanForBookExample(t *testing.T) {
	result := bookResult(t)
	if want := 184.0 / 35; !near(result.Mean, want, exact) {
		t.Errorf("E*T = %v, want %v", result.Mean, want)
	}
	// Trapezoidal rule: (35 + 2·149 + 0) / (2·35).
	if want := 333.0 / 70; !near(result.MeanFromReliability, want, exact) {
		t.Errorf("E*T from R*(t) = %v, want %v", result.MeanFromReliability, want)
	}
	if result.Failed != bookSampleSize || result.Survived != 0 {
		t.Errorf("failed %d, survived %d; want %d and 0", result.Failed, result.Survived, bookSampleSize)
	}
}

func TestComputeHandlesUnevenIntervals(t *testing.T) {
	result, err := Compute(10, []Observation{{2, 2}, {5, 3}, {6, 5}})
	if err != nil {
		t.Fatal(err)
	}
	second := result.Rows[1]
	// f*(t_2) = 3 / (10·3), λ*(t_2) = 3 / (8·3)
	if !near(second.Density, 0.1, exact) || !near(second.Hazard, 0.125, exact) {
		t.Errorf("f* = %v, λ* = %v; want 0.1 and 0.125", second.Density, second.Hazard)
	}
	if want := (2*2 + 5*3 + 6*5) / 10.0; !near(result.Mean, want, exact) {
		t.Errorf("E*T = %v, want %v", result.Mean, want)
	}
}

func TestComputeHandlesSurvivingElements(t *testing.T) {
	result, err := Compute(10, []Observation{{1, 2}, {2, 3}})
	if err != nil {
		t.Fatal(err)
	}
	if result.Failed != 5 || result.Survived != 5 {
		t.Errorf("failed %d, survived %d; want 5 and 5", result.Failed, result.Survived)
	}
	if !math.IsNaN(result.Mean) {
		t.Errorf("E*T = %v, want undefined", result.Mean)
	}
	// (1 + 0.8)/2 · 1 + (0.8 + 0.5)/2 · 1
	if !near(result.MeanFromReliability, 1.55, exact) {
		t.Errorf("E*T from R*(t) = %v, want 1.55", result.MeanFromReliability)
	}
	if !near(result.Rows[1].Reliability, 0.5, exact) {
		t.Errorf("R*(t_2) = %v, want 0.5", result.Rows[1].Reliability)
	}
}

func TestComputeRejectsInvalidData(t *testing.T) {
	cases := []struct {
		name         string
		sampleSize   int
		observations []Observation
		want         error
	}{
		{"zero sample size", 0, bookObservations, ErrSampleSizeNotPositive},
		{"no observations", 5, nil, ErrNoObservations},
		{"repeated time", 5, []Observation{{2, 1}, {2, 1}}, ErrTimeNotIncreasing},
		{"zero time", 5, []Observation{{0, 1}}, ErrTimeNotIncreasing},
		{"NaN time", 5, []Observation{{math.NaN(), 1}}, ErrTimeNotIncreasing},
		{"negative failures", 5, []Observation{{1, -1}}, ErrNegativeFailures},
		{"too many failures", 5, []Observation{{1, 3}, {2, 3}}, ErrFailuresExceedSample},
	}
	for _, c := range cases {
		if _, err := Compute(c.sampleSize, c.observations); !errors.Is(err, c.want) {
			t.Errorf("%s: error = %v, want %v", c.name, err, c.want)
		}
	}
}

func TestComputeReportsPositionOfInvalidObservation(t *testing.T) {
	_, err := Compute(5, []Observation{{1, 1}, {2, 1}, {2, 1}})
	var item *ItemError
	if !errors.As(err, &item) || item.Position != 3 {
		t.Errorf("error = %v, want an ItemError at position 3", err)
	}
}

func TestGroupCountsTimesPerInterval(t *testing.T) {
	observations, err := Group([]float64{0.5, 1, 1.2, 2, 3.7}, 1)
	if err != nil {
		t.Fatal(err)
	}
	want := []Observation{{1, 2}, {2, 2}, {3, 0}, {4, 1}}
	if len(observations) != len(want) {
		t.Fatalf("got %d intervals, want %d", len(observations), len(want))
	}
	for i := range want {
		if observations[i] != want[i] {
			t.Errorf("interval %d = %+v, want %+v", i+1, observations[i], want[i])
		}
	}
}

func TestGroupKeepsBoundaryTimeInLowerInterval(t *testing.T) {
	// 0.3 lies on the boundary 3·0.1 even though 0.3/0.1 rounds above 3.
	observations, err := Group([]float64{0.1, 0.2, 0.3}, 0.1)
	if err != nil {
		t.Fatal(err)
	}
	want := []int{1, 1, 1}
	if len(observations) != len(want) {
		t.Fatalf("got %+v, want %d intervals", observations, len(want))
	}
	for i, failures := range want {
		if observations[i].Failures != failures {
			t.Errorf("interval %d holds %d failures, want %d", i+1, observations[i].Failures, failures)
		}
	}
}

func TestGroupRejectsInvalidData(t *testing.T) {
	cases := []struct {
		name  string
		times []float64
		step  float64
		want  error
	}{
		{"no times", nil, 1, ErrNoFailureTimes},
		{"zero step", []float64{1, 2}, 0, ErrStepNotPositive},
		{"NaN step", []float64{1, 2}, math.NaN(), ErrStepNotPositive},
		{"negative time", []float64{1, -2}, 1, ErrFailureTimeNotPositive},
		{"too many intervals", []float64{1e9}, 1, ErrTooManyIntervals},
	}
	for _, c := range cases {
		if _, err := Group(c.times, c.step); !errors.Is(err, c.want) {
			t.Errorf("%s: error = %v, want %v", c.name, err, c.want)
		}
	}
}

func TestFromTimesUsesExactMean(t *testing.T) {
	times := []float64{4, 1.5, 2.5}
	result, err := FromTimes(times, len(times), FixedStep(1))
	if err != nil {
		t.Fatal(err)
	}
	if want := 8.0 / 3; !near(result.Mean, want, exact) {
		t.Errorf("E*T = %v, want %v", result.Mean, want)
	}
	if result.SampleSize != 3 || result.Survived != 0 {
		t.Errorf("sample size %d, survived %d; want 3 and 0", result.SampleSize, result.Survived)
	}
	if result.Times[0] != 1.5 || result.Times[2] != 4 {
		t.Errorf("times = %v, want ascending order", result.Times)
	}
	if times[0] != 4 {
		t.Errorf("input slice was reordered: %v", times)
	}
}

func TestFromTimesHandlesSurvivingElements(t *testing.T) {
	result, err := FromTimes([]float64{1.5, 2.5, 4}, 5, FixedStep(1))
	if err != nil {
		t.Fatal(err)
	}
	if result.Survived != 2 || !math.IsNaN(result.Mean) {
		t.Errorf("survived %d, E*T = %v; want 2 and undefined", result.Survived, result.Mean)
	}
}

func TestFromTimesRejectsInvalidData(t *testing.T) {
	if _, err := FromTimes([]float64{1, 2, 3}, 2, AutoStep); !errors.Is(err, ErrSampleSmallerThanTimes) {
		t.Errorf("small sample: error = %v, want %v", err, ErrSampleSmallerThanTimes)
	}
	if _, err := FromTimes(nil, 0, AutoStep); !errors.Is(err, ErrNoFailureTimes) {
		t.Errorf("no times: error = %v, want %v", err, ErrNoFailureTimes)
	}
	if _, err := FromTimes([]float64{1, 2}, 2, FixedStep(-1)); !errors.Is(err, ErrStepNotPositive) {
		t.Errorf("negative step: error = %v, want %v", err, ErrStepNotPositive)
	}
}

func TestAutoStepRoundsUpToRoundNumber(t *testing.T) {
	cases := []struct {
		name    string
		count   int
		longest float64
		want    float64
	}{
		{"35 elements: 6 intervals, 9/6 = 1.5", 35, 9, 2},
		{"4 elements: minimum of 5 intervals, 10/5", 4, 10, 2},
		{"10000 elements: maximum of 50 intervals, 600/50", 10000, 600, 20},
		{"fractional length: 3/10 = 0.3", 100, 3, 0.5},
	}
	for _, c := range cases {
		times := make([]float64, c.count)
		for i := range times {
			times[i] = c.longest
		}
		if got := AutoStep(times); !near(got, c.want, exact) {
			t.Errorf("%s: step = %v, want %v", c.name, got, c.want)
		}
	}
}
