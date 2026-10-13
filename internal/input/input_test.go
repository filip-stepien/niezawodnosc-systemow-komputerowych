package input

import (
	"errors"
	"slices"
	"testing"

	"empirical-reliability/internal/numfmt"
	"empirical-reliability/internal/reliability"
)

func TestParseTimesAcceptsEverySeparator(t *testing.T) {
	got, err := ParseTimes("12,5 30;41.2\n55\t60")
	if err != nil {
		t.Fatal(err)
	}
	if want := []float64{12.5, 30, 41.2, 55, 60}; !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestParseTimesReturnsNothingForBlankText(t *testing.T) {
	got, err := ParseTimes("  \n ")
	if err != nil || len(got) != 0 {
		t.Errorf("got %v, %v; want no times and no error", got, err)
	}
}

func TestParseTimesReportsPositionOfMalformedNumber(t *testing.T) {
	_, err := ParseTimes("1 x 3")
	var item *reliability.ItemError
	if !errors.As(err, &item) || item.Position != 2 || !errors.Is(err, numfmt.ErrNotNumber) {
		t.Errorf("error = %v, want %v at position 2", err, numfmt.ErrNotNumber)
	}
}

func TestParseTableSkipsBlankLines(t *testing.T) {
	got, err := ParseTable("1;0\r\n\r\n2,5;3\r\n")
	if err != nil {
		t.Fatal(err)
	}
	want := []reliability.Observation{{Time: 1, Failures: 0}, {Time: 2.5, Failures: 3}}
	if !slices.Equal(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestParseTableRejectsMalformedRows(t *testing.T) {
	cases := []struct {
		name string
		text string
		want error
		line int
	}{
		{"one column", "1 2\n3", ErrColumnCount, 2},
		{"three columns", "1 2 3", ErrColumnCount, 1},
		{"malformed time", "a 2", numfmt.ErrNotNumber, 1},
		{"fractional count", "1 2\n\n2 2,5", numfmt.ErrNotInteger, 3},
	}
	for _, c := range cases {
		_, err := ParseTable(c.text)
		var line *LineError
		if !errors.As(err, &line) || line.Line != c.line || !errors.Is(err, c.want) {
			t.Errorf("%s: error = %v, want %v on line %d", c.name, err, c.want, c.line)
		}
	}
}

func TestFormatTableRoundTrips(t *testing.T) {
	want := []reliability.Observation{{Time: 0.5, Failures: 0}, {Time: 2, Failures: 3}}
	text := FormatTable(want)
	if text != "0,5 0\n2 3" {
		t.Errorf("FormatTable = %q", text)
	}
	got, err := ParseTable(text)
	if err != nil || !slices.Equal(got, want) {
		t.Errorf("round trip gave %+v, %v; want %+v", got, err, want)
	}
}
