// Package numfmt parses and formats numbers the way the user interface shows
// them: with a decimal comma. A decimal point is accepted on input as well.
package numfmt

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Undefined is shown in place of a value that cannot be determined.
const Undefined = "—"

const (
	displaySeparator = ","
	parseSeparator   = "."
	compactDigits    = 6
	floatBits        = 64
	integerBase      = 10
)

var (
	ErrNotNumber   = errors.New("not a number")
	ErrNotInteger  = errors.New("not an integer")
	ErrNotUnsigned = errors.New("not a non-negative integer")
)

// SyntaxError reports text that could not be read as a number.
type SyntaxError struct {
	Text string
	Err  error
}

func (e *SyntaxError) Error() string { return fmt.Sprintf("%q: %v", e.Text, e.Err) }
func (e *SyntaxError) Unwrap() error { return e.Err }

// ParseFloat reads a finite number written with a decimal comma or point.
func ParseFloat(text string) (float64, error) {
	text = strings.TrimSpace(text)
	value, err := strconv.ParseFloat(strings.ReplaceAll(text, displaySeparator, parseSeparator), floatBits)
	if err != nil || math.IsNaN(value) || math.IsInf(value, 0) {
		return 0, &SyntaxError{Text: text, Err: ErrNotNumber}
	}
	return value, nil
}

// ParseInt reads an integer.
func ParseInt(text string) (int, error) {
	text = strings.TrimSpace(text)
	value, err := strconv.Atoi(text)
	if err != nil {
		return 0, &SyntaxError{Text: text, Err: ErrNotInteger}
	}
	return value, nil
}

// ParseUint reads a non-negative integer.
func ParseUint(text string) (uint64, error) {
	text = strings.TrimSpace(text)
	value, err := strconv.ParseUint(text, integerBase, floatBits)
	if err != nil {
		return 0, &SyntaxError{Text: text, Err: ErrNotUnsigned}
	}
	return value, nil
}

// Fixed formats value with the given number of decimal places.
func Fixed(value float64, decimals int) string {
	if !isFinite(value) {
		return Undefined
	}
	return Localize(strconv.FormatFloat(value, 'f', decimals, floatBits))
}

// Compact formats value without trailing zeros, using at most six significant
// digits.
func Compact(value float64) string {
	if !isFinite(value) {
		return Undefined
	}
	return Localize(strconv.FormatFloat(value, 'g', compactDigits, floatBits))
}

// Localize replaces the decimal point of an already formatted number with the
// decimal comma.
func Localize(number string) string {
	return strings.ReplaceAll(number, parseSeparator, displaySeparator)
}

func isFinite(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0)
}
