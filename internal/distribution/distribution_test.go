package distribution

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

const (
	largeSample   = 50000
	meanTolerance = 0.02 // relative
)

func newRNG() *rand.Rand { return rand.New(rand.NewPCG(1, 2)) }

func mustBuild(t *testing.T, spec Spec) Distribution {
	t.Helper()
	values := make([]float64, len(spec.Parameters))
	for i, p := range spec.Parameters {
		values[i] = p.Default
	}
	d, err := spec.New(values...)
	if err != nil {
		t.Fatalf("%s: %v", spec.ID, err)
	}
	return d
}

func TestSpecRejectsNonPositiveParameters(t *testing.T) {
	cases := []struct {
		name   string
		spec   Spec
		values []float64
		symbol string
	}{
		{"zero rate", Exponential, []float64{0}, "λ"},
		{"negative shape", Weibull, []float64{-1, 10}, "β"},
		{"zero scale", Weibull, []float64{2, 0}, "η"},
		{"negative mean", Normal, []float64{-5, 1}, "μ"},
		{"zero deviation", Normal, []float64{100, 0}, "σ"},
		{"NaN rate", Exponential, []float64{math.NaN()}, "λ"},
		{"infinite rate", Exponential, []float64{math.Inf(1)}, "λ"},
	}
	for _, c := range cases {
		_, err := c.spec.New(c.values...)
		var parameter *ParameterError
		if !errors.As(err, &parameter) || parameter.Symbol != c.symbol {
			t.Errorf("%s: error = %v, want a ParameterError for %s", c.name, err, c.symbol)
		}
	}
}

func TestSpecRejectsWrongParameterCount(t *testing.T) {
	if _, err := Weibull.New(2); !errors.Is(err, ErrParameterCount) {
		t.Errorf("error = %v, want %v", err, ErrParameterCount)
	}
}

func TestCatalogHasUniqueIDs(t *testing.T) {
	seen := map[string]bool{}
	for _, spec := range Catalog() {
		if seen[spec.ID] {
			t.Errorf("duplicate distribution ID %q", spec.ID)
		}
		seen[spec.ID] = true
	}
}

func TestSampleMeanConvergesToTheoreticalMean(t *testing.T) {
	for _, spec := range Catalog() {
		d := mustBuild(t, spec)
		times, err := Sample(d, largeSample, newRNG())
		if err != nil {
			t.Fatal(err)
		}
		total := 0.0
		for _, v := range times {
			if v <= 0 {
				t.Fatalf("%s: drew a non-positive lifetime %v", spec.ID, v)
			}
			total += v
		}
		got := total / float64(len(times))
		if relative := math.Abs(got-d.Mean()) / d.Mean(); relative > meanTolerance {
			t.Errorf("%s: sample mean %.3f, theoretical mean %.3f", spec.ID, got, d.Mean())
		}
	}
}

func TestCDFIsMonotonicFromZeroToOne(t *testing.T) {
	for _, spec := range Catalog() {
		d := mustBuild(t, spec)
		if got := d.CDF(0); got != 0 {
			t.Errorf("%s: F(0) = %v, want 0", spec.ID, got)
		}
		previous := 0.0
		for x := 1.0; x <= 2000; x += 1 {
			current := d.CDF(x)
			if current < previous || current > 1 {
				t.Fatalf("%s: F(%v) = %v after %v", spec.ID, x, current, previous)
			}
			previous = current
		}
		if previous < 0.999 {
			t.Errorf("%s: F(2000) = %v, want close to 1", spec.ID, previous)
		}
	}
}

func TestSampleRejectsInvalidSize(t *testing.T) {
	d := mustBuild(t, Exponential)
	if _, err := Sample(d, 0, newRNG()); !errors.Is(err, ErrSampleSizeNotPositive) {
		t.Errorf("zero size: error = %v, want %v", err, ErrSampleSizeNotPositive)
	}
	if _, err := Sample(d, MaxSampleSize+1, newRNG()); !errors.Is(err, ErrSampleSizeTooLarge) {
		t.Errorf("oversized: error = %v, want %v", err, ErrSampleSizeTooLarge)
	}
}

func TestHazardOfExponentialIsConstant(t *testing.T) {
	const rate = 0.25
	d, err := Exponential.New(rate)
	if err != nil {
		t.Fatal(err)
	}
	for _, x := range []float64{0.1, 1, 10} {
		if got := Hazard(d, x); math.Abs(got-rate) > 1e-9 {
			t.Errorf("λ(%v) = %v, want %v", x, got, rate)
		}
	}
}

func TestHazardIsUndefinedWhereReliabilityVanishes(t *testing.T) {
	d := mustBuild(t, Exponential)
	if got := Hazard(d, 1e6); !math.IsNaN(got) {
		t.Errorf("λ(1e6) = %v, want undefined", got)
	}
}
