package ui

import (
	"strconv"

	"empirical-reliability/internal/numfmt"
	"empirical-reliability/internal/reliability"
)

// resultColumn describes one column of the results table. Showing another
// value only takes a new entry in resultColumns.
type resultColumn struct {
	header string
	width  float32
	cell   func(index int, row reliability.Row) string
}

var resultColumns = []resultColumn{
	{"i", 50, func(index int, _ reliability.Row) string { return strconv.Itoa(index + 1) }},
	{"t_i", 90, func(_ int, row reliability.Row) string { return numfmt.Compact(row.Time) }},
	{"Δt", 80, func(_ int, row reliability.Row) string { return numfmt.Compact(row.Step) }},
	{"m(t_i)", 70, func(_ int, row reliability.Row) string { return strconv.Itoa(row.Failures) }},
	{"n(t_i)", 70, func(_ int, row reliability.Row) string { return strconv.Itoa(row.Working) }},
	{"F*(t_i)", 95, func(_ int, row reliability.Row) string { return numfmt.Fixed(row.CDF, probabilityDecimals) }},
	{"R*(t_i)", 95, func(_ int, row reliability.Row) string { return numfmt.Fixed(row.Reliability, probabilityDecimals) }},
	{"f*(t_i)", 105, func(_ int, row reliability.Row) string { return numfmt.Fixed(row.Density, rateDecimals) }},
	{"λ*(t_i)", 105, func(_ int, row reliability.Row) string { return numfmt.Fixed(row.Hazard, rateDecimals) }},
}

func resultHeaders() []string {
	headers := make([]string, len(resultColumns))
	for i, column := range resultColumns {
		headers[i] = column.header
	}
	return headers
}

func resultCells(index int, row reliability.Row) []string {
	cells := make([]string, len(resultColumns))
	for i, column := range resultColumns {
		cells[i] = column.cell(index, row)
	}
	return cells
}

// exportRecords returns the results table followed by the summary values, in
// the form written to a CSV file.
func exportRecords(result *reliability.Result) [][]string {
	records := [][]string{resultHeaders()}
	for i, row := range result.Rows {
		records = append(records, resultCells(i, row))
	}
	return append(records,
		[]string{},
		[]string{exportSampleSize, strconv.Itoa(result.SampleSize)},
		[]string{exportMean, numfmt.Fixed(result.Mean, exportMeanDecimals)},
		[]string{exportMeanFromReliability, numfmt.Fixed(result.MeanFromReliability, exportMeanDecimals)},
	)
}
