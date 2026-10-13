package numfmt

import (
	"errors"
	"math"
	"testing"
)

func TestParseFloatAcceptsCommaAndPoint(t *testing.T) {
	cases := map[string]float64{"0,25": 0.25, "0.25": 0.25, " 12 ": 12, "1e2": 100, "-3,5": -3.5}
	for text, want := range cases {
		got, err := ParseFloat(text)
		if err != nil || got != want {
			t.Errorf("ParseFloat(%q) = %v, %v; want %v", text, got, err, want)
		}
	}
}

func TestParseFloatRejectsMalformedText(t *testing.T) {
	for _, text := range []string{"", "abc", "1,2,3", "NaN", "Inf", "1 2"} {
		if _, err := ParseFloat(text); !errors.Is(err, ErrNotNumber) {
			t.Errorf("ParseFloat(%q): error = %v, want %v", text, err, ErrNotNumber)
		}
	}
}

func TestParseIntRejectsFractions(t *testing.T) {
	if got, err := ParseInt(" 35 "); err != nil || got != 35 {
		t.Errorf("ParseInt = %v, %v; want 35", got, err)
	}
	for _, text := range []string{"", "2,5", "x"} {
		if _, err := ParseInt(text); !errors.Is(err, ErrNotInteger) {
			t.Errorf("ParseInt(%q): error = %v, want %v", text, err, ErrNotInteger)
		}
	}
}

func TestParseUintRejectsNegativeNumbers(t *testing.T) {
	if got, err := ParseUint("7"); err != nil || got != 7 {
		t.Errorf("ParseUint = %v, %v; want 7", got, err)
	}
	if _, err := ParseUint("-1"); !errors.Is(err, ErrNotUnsigned) {
		t.Errorf("error = %v, want %v", err, ErrNotUnsigned)
	}
}

func TestSyntaxErrorCarriesOffendingText(t *testing.T) {
	_, err := ParseFloat(" abc ")
	var syntax *SyntaxError
	if !errors.As(err, &syntax) || syntax.Text != "abc" {
		t.Errorf("error = %v, want a SyntaxError for %q", err, "abc")
	}
}

func TestFixedUsesDecimalComma(t *testing.T) {
	if got := Fixed(0.08571, 3); got != "0,086" {
		t.Errorf("Fixed = %q, want %q", got, "0,086")
	}
	if got := Fixed(math.NaN(), 3); got != Undefined {
		t.Errorf("Fixed(NaN) = %q, want %q", got, Undefined)
	}
}

func TestCompactDropsTrailingZeros(t *testing.T) {
	cases := map[float64]string{2.5: "2,5", 10: "10", 0.30000000000000004: "0,3"}
	for value, want := range cases {
		if got := Compact(value); got != want {
			t.Errorf("Compact(%v) = %q, want %q", value, got, want)
		}
	}
	if got := Compact(math.Inf(1)); got != Undefined {
		t.Errorf("Compact(Inf) = %q, want %q", got, Undefined)
	}
}
