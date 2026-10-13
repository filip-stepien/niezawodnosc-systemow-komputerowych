// Package input converts text typed by the user into data for the analysis.
package input

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"empirical-reliability/internal/numfmt"
	"empirical-reliability/internal/reliability"
)

const (
	fieldSeparators = " \t\r\n;"
	lineSeparator   = "\n"
	columnSeparator = " "
	tableColumns    = 2
)

// ErrColumnCount reports a table row that does not hold exactly a time and a
// failure count.
var ErrColumnCount = errors.New("a table row must hold a time and a failure count")

// LineError locates an error on a 1-based line of the failure table.
type LineError struct {
	Line int
	Err  error
}

func (e *LineError) Error() string { return fmt.Sprintf("line %d: %v", e.Line, e.Err) }
func (e *LineError) Unwrap() error { return e.Err }

// ParseTimes reads failure times separated by whitespace or semicolons.
func ParseTimes(text string) ([]float64, error) {
	fields := splitFields(text)
	times := make([]float64, len(fields))
	for i, field := range fields {
		value, err := numfmt.ParseFloat(field)
		if err != nil {
			return nil, &reliability.ItemError{Position: i + 1, Err: err}
		}
		times[i] = value
	}
	return times, nil
}

// ParseTable reads a failure table: every non-empty line holds an inspection
// time and the number of failures found at it.
func ParseTable(text string) ([]reliability.Observation, error) {
	var observations []reliability.Observation
	for i, line := range strings.Split(text, lineSeparator) {
		fields := splitFields(line)
		if len(fields) == 0 {
			continue
		}
		observation, err := parseRow(fields)
		if err != nil {
			return nil, &LineError{Line: i + 1, Err: err}
		}
		observations = append(observations, observation)
	}
	return observations, nil
}

// FormatTable writes observations in the form ParseTable reads.
func FormatTable(observations []reliability.Observation) string {
	lines := make([]string, len(observations))
	for i, o := range observations {
		lines[i] = numfmt.Compact(o.Time) + columnSeparator + strconv.Itoa(o.Failures)
	}
	return strings.Join(lines, lineSeparator)
}

func parseRow(fields []string) (reliability.Observation, error) {
	if len(fields) != tableColumns {
		return reliability.Observation{}, ErrColumnCount
	}
	time, err := numfmt.ParseFloat(fields[0])
	if err != nil {
		return reliability.Observation{}, err
	}
	failures, err := numfmt.ParseInt(fields[1])
	if err != nil {
		return reliability.Observation{}, err
	}
	return reliability.Observation{Time: time, Failures: failures}, nil
}

func splitFields(text string) []string {
	return strings.FieldsFunc(text, func(r rune) bool {
		return strings.ContainsRune(fieldSeparators, r)
	})
}
