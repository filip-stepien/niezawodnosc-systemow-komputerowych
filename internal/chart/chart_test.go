package chart

import (
	"math"
	"math/rand/v2"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"empirical-reliability/internal/distribution"
	"empirical-reliability/internal/export"
	"empirical-reliability/internal/reliability"
)

// outputDirVariable names the environment variable that, when set, makes the
// tests save the rendered charts to the given directory for inspection.
const outputDirVariable = "CHART_OUTPUT_DIR"

var testLabels = Labels{Title: "Title", TimeAxis: "t", Empirical: "empirical", Theoretical: "theoretical"}

func bookResult(t *testing.T) *reliability.Result {
	t.Helper()
	result, err := reliability.Compute(35, []reliability.Observation{
		{Time: 1, Failures: 0}, {Time: 2, Failures: 3}, {Time: 3, Failures: 3}, {Time: 4, Failures: 5},
		{Time: 5, Failures: 8}, {Time: 6, Failures: 7}, {Time: 7, Failures: 6}, {Time: 8, Failures: 2},
		{Time: 9, Failures: 1}, {Time: 10, Failures: 0},
	})
	if err != nil {
		t.Fatal(err)
	}
	return result
}

func randomResult(t *testing.T) (*reliability.Result, distribution.Distribution) {
	t.Helper()
	model, err := distribution.Weibull.New(2, 100)
	if err != nil {
		t.Fatal(err)
	}
	times, err := distribution.Sample(model, 500, rand.New(rand.NewPCG(1, 2)))
	if err != nil {
		t.Fatal(err)
	}
	result, err := reliability.FromTimes(times, len(times), reliability.AutoStep)
	if err != nil {
		t.Fatal(err)
	}
	return result, model
}

func TestRenderDrawsEveryIndicator(t *testing.T) {
	random, model := randomResult(t)
	cases := []struct {
		name   string
		result *reliability.Result
		model  distribution.Model
	}{
		{"book", bookResult(t), nil},
		{"weibull", random, model},
	}
	wantWidth, wantHeight := int(imageWidth.Dots(resolution)), int(imageHeight.Dots(resolution))

	for _, c := range cases {
		for _, indicator := range Indicators() {
			img, err := Render(indicator, c.result, c.model, testLabels)
			if err != nil {
				t.Fatalf("%s, %s: %v", c.name, indicator.Slug(), err)
			}
			if size := img.Bounds().Size(); size.X != wantWidth || size.Y != wantHeight {
				t.Errorf("%s, %s: image size %v, want %dx%d", c.name, indicator.Slug(), size, wantWidth, wantHeight)
			}
			if dir := os.Getenv(outputDirVariable); dir != "" {
				path := filepath.Join(dir, c.name+"_"+indicator.Slug()+".png")
				if err := export.SavePNG(path, img); err != nil {
					t.Fatal(err)
				}
			}
		}
	}
}

func TestPointsStopAtFirstUndefinedHazard(t *testing.T) {
	rows := bookResult(t).Rows
	points := definitions[Hazard].points(rows)
	// The origin plus nine defined values; λ*(t_10) is undefined.
	if len(points) != len(rows) {
		t.Fatalf("got %d points, want %d", len(points), len(rows))
	}
	for _, point := range points {
		if math.IsNaN(point.Y) {
			t.Errorf("undefined value at t = %v was not dropped", point.X)
		}
	}
}

func TestPointsStartAtTheValueForTimeZero(t *testing.T) {
	rows := bookResult(t).Rows
	if start := definitions[CDF].points(rows)[0]; start.X != 0 || start.Y != 0 {
		t.Errorf("F* starts at %+v, want (0, 0)", start)
	}
	if start := definitions[Reliability].points(rows)[0]; start.X != 0 || start.Y != 1 {
		t.Errorf("R* starts at %+v, want (0, 1)", start)
	}
	if start := definitions[Density].points(rows)[0]; start.Y != rows[0].Density {
		t.Errorf("f* starts at %v, want the first interval's value %v", start.Y, rows[0].Density)
	}
}

func TestUpperBound(t *testing.T) {
	rows := bookResult(t).Rows
	if got := definitions[CDF].upperBound(definitions[CDF].points(rows)); got != probabilityMax {
		t.Errorf("F* upper bound = %v, want %v", got, probabilityMax)
	}
	if got := definitions[Hazard].upperBound(definitions[Hazard].points(rows)); got != 1 {
		t.Errorf("λ* upper bound = %v, want the highest value 1", got)
	}
	flat, err := reliability.Compute(5, []reliability.Observation{{Time: 1, Failures: 0}})
	if err != nil {
		t.Fatal(err)
	}
	if got := definitions[Density].upperBound(definitions[Density].points(flat.Rows)); got != emptyRangeMax {
		t.Errorf("upper bound of an all-zero indicator = %v, want %v", got, emptyRangeMax)
	}
}

func TestTicksFallOnRoundValues(t *testing.T) {
	cases := []struct {
		name      string
		low, high float64
		want      []string
	}{
		{"time axis", 0, 260, []string{"0", "50", "100", "150", "200", "250"}},
		{"probability axis", 0, 1.08, []string{"0", "0,2", "0,4", "0,6", "0,8", "1"}},
		{"small rates", 0, 0.0108, []string{"0", "0,002", "0,004", "0,006", "0,008", "0,01"}},
		{"tick on the end of the axis", 0, 10, []string{"0", "2", "4", "6", "8", "10"}},
	}
	for _, c := range cases {
		ticks := roundTicks{}.Ticks(c.low, c.high)
		labels := make([]string, len(ticks))
		for i, tick := range ticks {
			labels[i] = tick.Label
		}
		if !slices.Equal(labels, c.want) {
			t.Errorf("%s: ticks %v, want %v", c.name, labels, c.want)
		}
	}
}

func TestEveryIndicatorHasUniqueSlug(t *testing.T) {
	seen := map[string]bool{}
	for _, indicator := range Indicators() {
		slug := indicator.Slug()
		if slug == "" || seen[slug] {
			t.Errorf("indicator %d has an empty or duplicate slug %q", indicator, slug)
		}
		seen[slug] = true
	}
}
