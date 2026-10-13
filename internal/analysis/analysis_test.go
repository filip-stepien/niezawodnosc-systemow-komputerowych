package analysis

import (
	"errors"
	"math"
	"slices"
	"testing"

	"empirical-reliability/internal/distribution"
	"empirical-reliability/internal/reliability"
)

const (
	largeSample  = 50000
	cdfTolerance = 0.01
)

func defaults(spec distribution.Spec) []float64 {
	values := make([]float64, len(spec.Parameters))
	for i, p := range spec.Parameters {
		values[i] = p.Default
	}
	return values
}

func request(spec distribution.Spec, size int, seed uint64) RandomRequest {
	return RandomRequest{Spec: spec, Parameters: defaults(spec), Size: size, Seed: seed, Step: reliability.AutoStep}
}

func TestFromTableHasNoTheoreticalModel(t *testing.T) {
	a, err := FromTable(5, []reliability.Observation{{Time: 1, Failures: 2}, {Time: 2, Failures: 3}})
	if err != nil {
		t.Fatal(err)
	}
	if a.Random != nil || a.Model() != nil {
		t.Errorf("manual data must not carry a model, got %+v", a.Random)
	}
	if a.Result.Failed != 5 {
		t.Errorf("failed = %d, want 5", a.Result.Failed)
	}
}

func TestFromTimesPropagatesValidationErrors(t *testing.T) {
	if _, err := FromTimes(nil, 0, reliability.AutoStep); !errors.Is(err, reliability.ErrNoFailureTimes) {
		t.Errorf("error = %v, want %v", err, reliability.ErrNoFailureTimes)
	}
}

func TestFromRandomIsReproducible(t *testing.T) {
	first, err := FromRandom(request(distribution.Weibull, 200, 7))
	if err != nil {
		t.Fatal(err)
	}
	second, err := FromRandom(request(distribution.Weibull, 200, 7))
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(first.Result.Times, second.Result.Times) {
		t.Error("the same seed produced different samples")
	}
	other, err := FromRandom(request(distribution.Weibull, 200, 8))
	if err != nil {
		t.Fatal(err)
	}
	if slices.Equal(first.Result.Times, other.Result.Times) {
		t.Error("different seeds produced the same sample")
	}
}

func TestFromRandomRecordsHowTheSampleWasDrawn(t *testing.T) {
	a, err := FromRandom(request(distribution.Normal, 50, 3))
	if err != nil {
		t.Fatal(err)
	}
	if a.Random.Spec.ID != distribution.Normal.ID || a.Random.Seed != 3 || a.Model() == nil {
		t.Errorf("random sample description = %+v", a.Random)
	}
	if a.Result.SampleSize != 50 || a.Result.Survived != 0 {
		t.Errorf("sample size %d, survived %d; want 50 and 0", a.Result.SampleSize, a.Result.Survived)
	}
}

func TestFromRandomRejectsInvalidRequests(t *testing.T) {
	invalid := request(distribution.Exponential, 100, 1)
	invalid.Parameters = []float64{-1}
	var parameter *distribution.ParameterError
	if _, err := FromRandom(invalid); !errors.As(err, &parameter) {
		t.Errorf("negative parameter: error = %v, want a ParameterError", err)
	}
	if _, err := FromRandom(request(distribution.Exponential, 0, 1)); !errors.Is(err, distribution.ErrSampleSizeNotPositive) {
		t.Errorf("zero size: error = %v, want %v", err, distribution.ErrSampleSizeNotPositive)
	}
}

// For a large sample the empirical indicators must approach the theoretical ones.
func TestFromRandomConvergesToTheory(t *testing.T) {
	for _, spec := range distribution.Catalog() {
		a, err := FromRandom(request(spec, largeSample, 1))
		if err != nil {
			t.Fatal(err)
		}
		model := a.Model()
		for _, row := range a.Result.Rows {
			if diff := math.Abs(row.CDF - model.CDF(row.Time)); diff > cdfTolerance {
				t.Errorf("%s: F*(%.1f) = %.4f, F = %.4f", spec.ID, row.Time, row.CDF, model.CDF(row.Time))
			}
		}
		want := a.Random.Model.Mean()
		if relative := math.Abs(a.Result.Mean-want) / want; relative > 0.02 {
			t.Errorf("%s: E*T = %.3f, ET = %.3f", spec.ID, a.Result.Mean, want)
		}
		if relative := math.Abs(a.Result.MeanFromReliability-want) / want; relative > 0.02 {
			t.Errorf("%s: E*T from R*(t) = %.3f, ET = %.3f", spec.ID, a.Result.MeanFromReliability, want)
		}
	}
}
