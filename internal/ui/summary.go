package ui

import (
	"fmt"
	"strings"

	"empirical-reliability/internal/analysis"
	"empirical-reliability/internal/numfmt"
	"empirical-reliability/internal/reliability"
)

const summaryLineSeparator = "\n"

// summaryText describes an analysis in a few lines shown above the results table.
func summaryText(a *analysis.Analysis) string {
	lines := sampleLines(a.Result)
	lines = append(lines, meanLines(a.Result)...)
	if a.Random != nil {
		lines = append(lines, randomLines(a.Random)...)
	}
	return strings.Join(lines, summaryLineSeparator)
}

func sampleLines(result *reliability.Result) []string {
	observationEnd := result.Rows[len(result.Rows)-1].Time
	return []string{
		fmt.Sprintf(summarySampleFormat, result.SampleSize, len(result.Rows), numfmt.Compact(observationEnd)),
		fmt.Sprintf(summaryCountsFormat, result.Failed, result.Survived),
	}
}

func meanLines(result *reliability.Result) []string {
	mean := numfmt.Fixed(result.Mean, meanDecimals)
	fromReliability := numfmt.Fixed(result.MeanFromReliability, meanDecimals)
	switch {
	case result.Survived > 0:
		return []string{summaryMeanUnknown, fmt.Sprintf(summaryMeanLowerBoundFormat, fromReliability)}
	case len(result.Times) == 0:
		// Grouped data carry no individual lifetimes.
		return []string{
			fmt.Sprintf(summaryMeanFromInspectionsFormat, mean),
			fmt.Sprintf(summaryMeanFromReliabilityFormat, fromReliability),
		}
	default:
		return []string{
			fmt.Sprintf(summaryMeanFormat, mean),
			fmt.Sprintf(summaryMeanFromReliabilityFormat, fromReliability),
		}
	}
}

func randomLines(sample *analysis.RandomSample) []string {
	parameters := make([]string, len(sample.Parameters))
	for i, value := range sample.Parameters {
		parameters[i] = fmt.Sprintf(summaryParameterFormat, sample.Spec.Parameters[i].Symbol, numfmt.Compact(value))
	}
	return []string{
		fmt.Sprintf(summaryDistributionFormat,
			distributionNames[sample.Spec.ID],
			strings.Join(parameters, summaryParameterSeparator),
			sample.Seed),
		fmt.Sprintf(summaryTheoreticalMeanFormat, numfmt.Fixed(sample.Model.Mean(), meanDecimals)),
	}
}
