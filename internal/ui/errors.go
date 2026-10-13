package ui

import (
	"errors"
	"fmt"

	"empirical-reliability/internal/distribution"
	"empirical-reliability/internal/input"
	"empirical-reliability/internal/numfmt"
	"empirical-reliability/internal/reliability"
)

const (
	errorFieldFormat     = "%s: %s"
	errorLineFormat      = "wiersz %d: %s"
	errorItemFormat      = "pozycja %d: %s"
	errorSyntaxFormat    = "„%s” %s"
	errorParameterFormat = "parametr %s musi być dodatni"
)

// errNothingToSave reports an attempt to save before anything was computed.
var errNothingToSave = errors.New("nothing to save")

// errorMessages translates the errors of the lower layers, which are written
// for developers, into messages for the user.
var errorMessages = map[error]string{
	reliability.ErrSampleSizeNotPositive:  "liczność próby n(0) musi być liczbą dodatnią",
	reliability.ErrNoObservations:         "tabela uszkodzeń jest pusta",
	reliability.ErrTimeNotIncreasing:      "chwile t muszą być dodatnie i rosnące",
	reliability.ErrNegativeFailures:       "liczba uszkodzeń nie może być ujemna",
	reliability.ErrFailuresExceedSample:   "łączna liczba uszkodzeń przekracza liczność próby n(0)",
	reliability.ErrNoFailureTimes:         "nie podano żadnych czasów uszkodzeń",
	reliability.ErrFailureTimeNotPositive: "czas uszkodzenia musi być dodatni",
	reliability.ErrStepNotPositive:        "długość przedziału Δt musi być dodatnia",
	reliability.ErrTooManyIntervals: fmt.Sprintf(
		"zbyt mała długość przedziału Δt (ponad %d przedziałów)", reliability.MaxIntervals),
	reliability.ErrSampleSmallerThanTimes: "liczność próby n(0) jest mniejsza od liczby podanych czasów",

	distribution.ErrSampleSizeNotPositive: "liczność próby musi być liczbą dodatnią",
	distribution.ErrSampleSizeTooLarge: fmt.Sprintf(
		"liczność próby nie może przekraczać %d", distribution.MaxSampleSize),
	distribution.ErrParameterCount: "niepoprawna liczba parametrów rozkładu",

	input.ErrColumnCount: "oczekiwano dwóch wartości: t oraz m(t)",

	numfmt.ErrNotNumber:   "nie jest poprawną liczbą",
	numfmt.ErrNotInteger:  "nie jest poprawną liczbą całkowitą",
	numfmt.ErrNotUnsigned: "nie jest nieujemną liczbą całkowitą",

	errNothingToSave: "najpierw wykonaj obliczenia",
}

// fieldError attaches the label of the form field an invalid value came from.
type fieldError struct {
	label string
	err   error
}

func (e *fieldError) Error() string { return fmt.Sprintf(errorFieldFormat, e.label, e.err) }
func (e *fieldError) Unwrap() error { return e.err }

// describe returns the message shown to the user for err. Errors that carry a
// location are unwrapped layer by layer, so the type switch deliberately
// inspects only the outermost error.
func describe(err error) string {
	switch e := err.(type) {
	case *fieldError:
		return fmt.Sprintf(errorFieldFormat, e.label, describe(e.err))
	case *input.LineError:
		return fmt.Sprintf(errorLineFormat, e.Line, describe(e.Err))
	case *reliability.ItemError:
		return fmt.Sprintf(errorItemFormat, e.Position, describe(e.Err))
	case *numfmt.SyntaxError:
		return fmt.Sprintf(errorSyntaxFormat, e.Text, describe(e.Err))
	case *distribution.ParameterError:
		return fmt.Sprintf(errorParameterFormat, e.Symbol)
	}
	for sentinel, message := range errorMessages {
		if errors.Is(err, sentinel) {
			return message
		}
	}
	return err.Error()
}
