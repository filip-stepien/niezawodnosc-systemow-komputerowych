package ui

import (
	"strings"

	"empirical-reliability/internal/analysis"
	"empirical-reliability/internal/numfmt"
	"empirical-reliability/internal/reliability"
)

// compute produces an analysis from the data currently entered in a form.
type compute func() (*analysis.Analysis, error)

// runner executes a computation and presents its outcome. Forms depend on
// this abstraction instead of on the window that implements it.
type runner func(compute)

// parseField reads the text of a form field, attaching the field's label to
// the error when the text is malformed.
func parseField[T any](label, text string, parse func(string) (T, error)) (T, error) {
	value, err := parse(text)
	if err != nil {
		return value, &fieldError{label: label, err: err}
	}
	return value, nil
}

func isBlank(text string) bool {
	return strings.TrimSpace(text) == ""
}

// parseStep returns the interval length entered by the user, or the automatic
// choice when the field is left blank.
func parseStep(text string) (reliability.StepSelector, error) {
	if isBlank(text) {
		return reliability.AutoStep, nil
	}
	step, err := parseField(labelStep, text, numfmt.ParseFloat)
	if err != nil {
		return nil, err
	}
	return reliability.FixedStep(step), nil
}

// parseSeed returns the seed entered by the user, or a fresh one when the
// field is left blank.
func parseSeed(text string) (uint64, error) {
	if isBlank(text) {
		return analysis.TimeSeed(), nil
	}
	return parseField(labelSeed, text, numfmt.ParseUint)
}
