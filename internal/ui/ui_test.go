package ui

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"

	"empirical-reliability/internal/analysis"
	"empirical-reliability/internal/chart"
	"empirical-reliability/internal/distribution"
	"empirical-reliability/internal/export"
	"empirical-reliability/internal/input"
	"empirical-reliability/internal/numfmt"
	"empirical-reliability/internal/reliability"
)

// screenshotDirVariable names the environment variable that, when set, makes
// the tests save screenshots of the window to the given directory.
const screenshotDirVariable = "SCREENSHOT_DIR"

// Positions of the tabs, in the order mainWindow adds them.
const (
	randomTab = iota
	tableTab
	timesTab
)

const (
	resultsTab = iota
	chartsTab
	sampleTab
)

// newTestWindow builds the main window on Fyne's in-memory test driver.
func newTestWindow(t *testing.T) (*mainWindow, fyne.Window) {
	t.Helper()
	application := test.NewApp()
	t.Cleanup(application.Quit)
	application.Settings().SetTheme(newLightTheme())

	window := application.NewWindow(windowTitle)
	m := newMainWindow(window)
	window.SetContent(m.content())
	window.Resize(fyne.NewSize(windowWidth, windowHeight))
	return m, window
}

func screenshot(t *testing.T, window fyne.Window, name string) {
	t.Helper()
	dir := os.Getenv(screenshotDirVariable)
	if dir == "" {
		return
	}
	if err := export.SavePNG(filepath.Join(dir, name), window.Canvas().Capture()); err != nil {
		t.Fatal(err)
	}
}

func TestTableFormComputesBookExample(t *testing.T) {
	m, window := newTestWindow(t)
	screenshot(t, window, "01_start.png")

	m.input.SelectIndex(tableTab)
	test.Tap(m.table.submit)

	if m.results.current == nil {
		t.Fatal("nothing was computed")
	}
	if !strings.Contains(m.results.summary.Text, "E*T ze wzoru (5-15) = 5,2571") {
		t.Errorf("summary = %q", m.results.summary.Text)
	}
	if rows, columns := m.results.table.Length(); rows != len(exampleObservations) || columns != len(resultColumns) {
		t.Errorf("results table is %d x %d", rows, columns)
	}
	if len(m.charts.rendered) != len(chart.Indicators()) {
		t.Errorf("%d charts rendered, want %d", len(m.charts.rendered), len(chart.Indicators()))
	}
	if m.sample.times.Text != "" {
		t.Errorf("grouped data must not list individual times, got %q", m.sample.times.Text)
	}
	screenshot(t, window, "02_table_results.png")
	m.output.SelectIndex(chartsTab)
	screenshot(t, window, "03_table_charts.png")
}

func TestRandomFormDrawsReproducibleSample(t *testing.T) {
	m, window := newTestWindow(t)
	m.random.seed.SetText("7")
	test.Tap(m.random.submit)

	first := m.results.current
	if first == nil || first.Random == nil {
		t.Fatal("no random sample was analysed")
	}
	if first.Result.SampleSize != 1000 || len(first.Result.Times) != 1000 {
		t.Errorf("sample size %d with %d times, want 1000", first.Result.SampleSize, len(first.Result.Times))
	}
	if !strings.Contains(m.results.summary.Text, "Rozkład wykładniczy, λ = 0,01; ziarno generatora: 7") {
		t.Errorf("summary = %q", m.results.summary.Text)
	}
	screenshot(t, window, "04_random_results.png")
	m.output.SelectIndex(chartsTab)
	screenshot(t, window, "05_random_charts.png")
	m.output.SelectIndex(sampleTab)
	screenshot(t, window, "06_random_sample.png")

	test.Tap(m.random.submit)
	if !slices.Equal(first.Result.Times, m.results.current.Result.Times) {
		t.Error("the same seed produced a different sample")
	}
}

func TestRandomFormShowsParametersOfSelectedDistribution(t *testing.T) {
	m, window := newTestWindow(t)
	m.random.seed.SetText("1")
	m.output.SelectIndex(chartsTab)
	for i, spec := range distribution.Catalog() {
		m.random.kind.SetSelectedIndex(i)
		if len(m.random.parameters) != len(spec.Parameters) {
			t.Fatalf("%s: %d parameter fields, want %d", spec.ID, len(m.random.parameters), len(spec.Parameters))
		}
		for j, parameter := range spec.Parameters {
			if got, want := m.random.parameters[j].Text, numfmt.Compact(parameter.Default); got != want {
				t.Errorf("%s: parameter %s defaults to %q, want %q", spec.ID, parameter.Symbol, got, want)
			}
		}
		a, err := m.random.analyze()
		if err != nil || a.Random.Spec.ID != spec.ID {
			t.Errorf("%s: analyze returned %v", spec.ID, err)
		}
		// The button must stay below the last field however many rows there are.
		lastField := m.random.seed.Position().Y + m.random.seed.Size().Height
		if button := m.random.submit.Position().Y; button < m.random.fields.Position().Y+lastField {
			t.Errorf("%s: the button at y = %v covers the fields ending at y = %v", spec.ID, button, lastField)
		}

		test.Tap(m.random.submit)
		if current := m.results.current; current == nil || current.Random.Spec.ID != spec.ID {
			t.Errorf("%s: the submitted sample was not presented", spec.ID)
		}
		screenshot(t, window, "07_random_"+spec.ID+"_charts.png")
	}
}

func TestTimesFormTreatsMissingElementsAsSurvivors(t *testing.T) {
	m, window := newTestWindow(t)
	m.input.SelectIndex(timesTab)
	m.times.size.SetText("12")
	m.times.step.SetText("10")
	m.times.data.SetText("12,5 30 41,2 55 60 61 70 85 90 120")
	test.Tap(m.times.submit)

	current := m.results.current
	if current == nil {
		t.Fatal("nothing was computed")
	}
	if current.Result.Survived != 2 || !strings.Contains(m.results.summary.Text, "oszacowanie od dołu") {
		t.Errorf("survived %d, summary %q", current.Result.Survived, m.results.summary.Text)
	}
	if lines := strings.Count(m.sample.times.Text, sampleSeparator) + 1; lines != 10 {
		t.Errorf("sample view lists %d times, want 10", lines)
	}
	screenshot(t, window, "08_times_results.png")
}

func TestTimesFormDefaultsSampleSizeToNumberOfTimes(t *testing.T) {
	m, _ := newTestWindow(t)
	m.times.data.SetText("1 2 3 4")
	a, err := m.times.analyze()
	if err != nil {
		t.Fatal(err)
	}
	if a.Result.SampleSize != 4 || a.Result.Survived != 0 {
		t.Errorf("sample size %d, survived %d; want 4 and 0", a.Result.SampleSize, a.Result.Survived)
	}
}

func TestTableFormClearAndExampleButtons(t *testing.T) {
	m, _ := newTestWindow(t)
	m.table.clear()
	if m.table.size.Text != "" || m.table.data.Text != "" {
		t.Error("clear left data in the form")
	}
	m.table.loadExample()
	if m.table.data.Text != input.FormatTable(exampleObservations) {
		t.Errorf("example data = %q", m.table.data.Text)
	}
}

func TestInvalidInputIsReportedInPolish(t *testing.T) {
	m, window := newTestWindow(t)
	cases := []struct {
		name    string
		prepare func()
		analyze compute
		want    string
	}{
		{
			"empty table",
			func() { m.table.data.SetText("") },
			m.table.analyze,
			"tabela uszkodzeń jest pusta",
		},
		{
			"malformed sample size",
			func() { m.table.loadExample(); m.table.size.SetText("abc") },
			m.table.analyze,
			"Liczność próby n(0): „abc” nie jest poprawną liczbą całkowitą",
		},
		{
			"more failures than elements",
			func() { m.table.loadExample(); m.table.size.SetText("10") },
			m.table.analyze,
			// The table holds 0, 3, 3 and 5 failures: the fourth row exceeds 10.
			"pozycja 4: łączna liczba uszkodzeń przekracza liczność próby n(0)",
		},
		{
			"malformed table row",
			func() { m.table.size.SetText("5"); m.table.data.SetText("1 2\n3") },
			m.table.analyze,
			"wiersz 2: oczekiwano dwóch wartości: t oraz m(t)",
		},
		{
			"negative distribution parameter",
			func() { m.random.parameters[0].SetText("-1") },
			m.random.analyze,
			"parametr λ musi być dodatni",
		},
		{
			"malformed seed",
			func() { m.random.kind.SetSelectedIndex(0); m.random.seed.SetText("x") },
			m.random.analyze,
			"Ziarno generatora: „x” nie jest nieujemną liczbą całkowitą",
		},
		{
			"zero random sample size",
			func() { m.random.seed.SetText(""); m.random.size.SetText("0") },
			m.random.analyze,
			"liczność próby musi być liczbą dodatnią",
		},
		{
			"zero step",
			func() { m.times.data.SetText("1 2 3"); m.times.step.SetText("0") },
			m.times.analyze,
			"długość przedziału Δt musi być dodatnia",
		},
		{
			"malformed time",
			func() { m.times.step.SetText(""); m.times.data.SetText("1 x 3") },
			m.times.analyze,
			"pozycja 2: „x” nie jest poprawną liczbą",
		},
		{
			"sample smaller than the list of times",
			func() { m.times.data.SetText("1 2 3"); m.times.size.SetText("2") },
			m.times.analyze,
			"liczność próby n(0) jest mniejsza od liczby podanych czasów",
		},
	}
	for _, c := range cases {
		c.prepare()
		_, err := c.analyze()
		if err == nil {
			t.Errorf("%s: want an error", c.name)
			continue
		}
		if got := describe(err); got != c.want {
			t.Errorf("%s: message = %q, want %q", c.name, got, c.want)
		}
	}

	m.input.SelectIndex(timesTab)
	test.Tap(m.times.submit)
	if window.Canvas().Overlays().Top() == nil {
		t.Error("an invalid form did not open an error dialog")
	}
	if m.results.current != nil {
		t.Error("an invalid form replaced the results")
	}
	screenshot(t, window, "09_error_dialog.png")
}

func TestEveryKnownErrorHasMessage(t *testing.T) {
	for sentinel, message := range errorMessages {
		if message == "" || describe(sentinel) != message {
			t.Errorf("%v is described as %q, want %q", sentinel, describe(sentinel), message)
		}
	}
	unknown := errors.New("disk failure")
	if got := describe(unknown); got != unknown.Error() {
		t.Errorf("unknown error described as %q", got)
	}
}

func TestSavingBeforeComputingIsRejected(t *testing.T) {
	m, window := newTestWindow(t)
	m.results.saveCSV()
	if window.Canvas().Overlays().Top() == nil {
		t.Error("saving results without data did not open an error dialog")
	}
}

func TestExportRecordsContainBookValues(t *testing.T) {
	a, err := analysis.FromTable(exampleSampleSize, exampleObservations)
	if err != nil {
		t.Fatal(err)
	}
	var buffer bytes.Buffer
	if err := export.WriteCSV(&buffer, exportRecords(a.Result)); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"i;t_i;Δt;m(t_i);n(t_i);F*(t_i);R*(t_i);f*(t_i);λ*(t_i)",
		"5;5;1;8;16;0,5429;0,4571;0,22857;0,33333",
		"10;10;1;0;0;1,0000;0,0000;0,00000;" + numfmt.Undefined,
		"n(0);35",
		"E*T (5-15);5,257143",
		"E*T z R*(t);4,757143",
	} {
		if !strings.Contains(buffer.String(), want) {
			t.Errorf("exported file lacks the line %q", want)
		}
	}
}

func TestChartsAreSavedUnderEnglishNames(t *testing.T) {
	m, _ := newTestWindow(t)
	test.Tap(m.table.submit)
	dir := t.TempDir()
	if err := m.charts.saveTo(dir); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"cdf.png", "reliability.png", "density.png", "hazard.png"} {
		if _, err := os.Stat(filepath.Join(dir, name)); err != nil {
			t.Errorf("chart file %s was not written: %v", name, err)
		}
	}
}

func TestSummaryDescribesEachKindOfData(t *testing.T) {
	grouped, err := analysis.FromTable(exampleSampleSize, exampleObservations)
	if err != nil {
		t.Fatal(err)
	}
	if text := summaryText(grouped); !strings.Contains(text, "przyjęto chwilę rejestracji uszkodzenia") {
		t.Errorf("grouped data summary = %q", text)
	}

	times, err := analysis.FromTimes([]float64{1, 2, 3}, 3, reliability.FixedStep(1))
	if err != nil {
		t.Fatal(err)
	}
	text := summaryText(times)
	if !strings.Contains(text, "E*T ze wzoru (5-15) = 2,0000\n") || strings.Contains(text, "Rozkład") {
		t.Errorf("individual times summary = %q", text)
	}
}
