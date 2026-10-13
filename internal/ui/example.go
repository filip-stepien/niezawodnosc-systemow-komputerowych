package ui

import "empirical-reliability/internal/reliability"

// The data of example 5.1 from the book (p. 209): 35 capacitors inspected
// every hour until the last one failed.
const exampleSampleSize = 35

var exampleObservations = []reliability.Observation{
	{Time: 1, Failures: 0},
	{Time: 2, Failures: 3},
	{Time: 3, Failures: 3},
	{Time: 4, Failures: 5},
	{Time: 5, Failures: 8},
	{Time: 6, Failures: 7},
	{Time: 7, Failures: 6},
	{Time: 8, Failures: 2},
	{Time: 9, Failures: 1},
	{Time: 10, Failures: 0},
}
