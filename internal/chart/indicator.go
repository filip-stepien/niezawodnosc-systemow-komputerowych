package chart

import (
	"math"

	"gonum.org/v1/plot/plotter"

	"empirical-reliability/internal/distribution"
	"empirical-reliability/internal/reliability"
)

// Indicator identifies the reliability indicator shown on a chart.
type Indicator int

const (
	CDF Indicator = iota
	Reliability
	Density
	Hazard
)

// Indicators lists the indicators that can be charted, in display order.
func Indicators() []Indicator {
	return []Indicator{CDF, Reliability, Density, Hazard}
}

// Slug returns a short identifier of the indicator, suitable for file names.
func (i Indicator) Slug() string { return definitions[i].slug }

// definition describes how one indicator is drawn. Charting a new indicator
// only takes a new entry in definitions.
type definition struct {
	slug   string
	value  func(reliability.Row) float64
	theory func(m distribution.Model, t float64) float64
	// interval marks an indicator that describes a whole interval
	// (t_{i-1}, t_i] and is therefore drawn as steps.
	interval bool
	// start is the value at t = 0 of an indicator defined at points.
	start      float64
	legendLeft bool
}

var definitions = map[Indicator]definition{
	CDF: {
		slug:       "cdf",
		value:      func(r reliability.Row) float64 { return r.CDF },
		theory:     func(m distribution.Model, t float64) float64 { return m.CDF(t) },
		start:      0,
		legendLeft: true,
	},
	Reliability: {
		slug:   "reliability",
		value:  func(r reliability.Row) float64 { return r.Reliability },
		theory: func(m distribution.Model, t float64) float64 { return 1 - m.CDF(t) },
		start:  1,
	},
	Density: {
		slug:     "density",
		value:    func(r reliability.Row) float64 { return r.Density },
		theory:   func(m distribution.Model, t float64) float64 { return m.PDF(t) },
		interval: true,
	},
	Hazard: {
		slug:       "hazard",
		value:      func(r reliability.Row) float64 { return r.Hazard },
		theory:     distribution.Hazard,
		interval:   true,
		legendLeft: true,
	},
}

// points returns the empirical values to draw, starting at t = 0 and stopping
// at the first undefined value.
func (d definition) points(rows []reliability.Row) plotter.XYs {
	start := d.start
	if d.interval {
		start = d.value(rows[0])
	}
	points := plotter.XYs{{X: 0, Y: start}}
	for _, row := range rows {
		value := d.value(row)
		if math.IsNaN(value) {
			break
		}
		points = append(points, plotter.XY{X: row.Time, Y: value})
	}
	return points
}

// upperBound returns the highest value the chart has to show.
func (d definition) upperBound(points plotter.XYs) float64 {
	if !d.interval {
		return probabilityMax
	}
	highest := 0.0
	for _, point := range points {
		highest = max(highest, point.Y)
	}
	if highest == 0 {
		return emptyRangeMax
	}
	return highest
}
