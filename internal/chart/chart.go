// Package chart draws the empirical reliability indicators, optionally
// together with the theoretical curve they estimate.
package chart

import (
	"image"
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"

	"empirical-reliability/internal/distribution"
	"empirical-reliability/internal/numfmt"
	"empirical-reliability/internal/reliability"
)

// Labels holds the texts shown on a chart.
type Labels struct {
	Title       string
	TimeAxis    string
	Empirical   string
	Theoretical string
}

// Render draws one indicator of result. When model is not nil, its curve is
// drawn as a dashed line for comparison.
func Render(indicator Indicator, result *reliability.Result, model distribution.Model, labels Labels) (image.Image, error) {
	definition := definitions[indicator]
	points := definition.points(result.Rows)
	lastTime := result.Rows[len(result.Rows)-1].Time

	p := newPlot(labels, lastTime, definition.upperBound(points)*headroom)
	line, err := addEmpirical(p, points, definition.interval)
	if err != nil {
		return nil, err
	}
	if model != nil {
		curve := addTheory(p, func(t float64) float64 { return definition.theory(model, t) })
		addLegend(p, labels, line, curve, definition.legendLeft)
	}
	return rasterize(p), nil
}

func newPlot(labels Labels, timeMax, valueMax float64) *plot.Plot {
	p := plot.New()
	p.Title.Text = labels.Title
	p.Title.Padding = titlePadding
	p.X.Label.Text = labels.TimeAxis
	p.X.Tick.Marker = roundTicks{}
	p.Y.Tick.Marker = roundTicks{}
	p.X.Min, p.X.Max = 0, timeMax
	p.Y.Min, p.Y.Max = 0, valueMax
	p.Add(plotter.NewGrid())
	return p
}

// addEmpirical draws the empirical values: as steps for interval indicators,
// otherwise as a line through marked points.
func addEmpirical(p *plot.Plot, points plotter.XYs, interval bool) (*plotter.Line, error) {
	line, err := plotter.NewLine(points)
	if err != nil {
		return nil, err
	}
	line.Color = empiricalColor
	line.Width = empiricalWidth
	p.Add(line)

	if interval {
		line.StepStyle = plotter.PreStep
		return line, nil
	}
	if len(points) <= maxMarkers {
		markers, err := plotter.NewScatter(points)
		if err != nil {
			return nil, err
		}
		markers.GlyphStyle = draw.GlyphStyle{Color: empiricalColor, Radius: markerRadius, Shape: draw.CircleGlyph{}}
		p.Add(markers)
	}
	return line, nil
}

func addTheory(p *plot.Plot, theory func(float64) float64) *plotter.Function {
	offChart := p.Y.Max * offChartFactor
	curve := plotter.NewFunction(func(t float64) float64 {
		value := theory(t)
		if math.IsNaN(value) || value > offChart {
			return offChart
		}
		return value
	})
	curve.Samples = theorySamples
	curve.Color = theoryColor
	curve.Width = theoryWidth
	curve.Dashes = theoryDashes
	p.Add(curve)
	return curve
}

func addLegend(p *plot.Plot, labels Labels, empirical, theoretical plot.Thumbnailer, left bool) {
	p.Legend.Add(labels.Empirical, empirical)
	p.Legend.Add(labels.Theoretical, theoretical)
	p.Legend.Top = true
	p.Legend.Left = left
	p.Legend.Padding = legendPadding
	p.Legend.TextStyle.Font.Size = legendFontSize
}

// rasterize draws the plot to an image. The margin keeps the outermost tick
// labels from being cut off.
func rasterize(p *plot.Plot) image.Image {
	canvas := vgimg.NewWith(vgimg.UseWH(imageWidth, imageHeight), vgimg.UseDPI(resolution))
	p.Draw(draw.Crop(draw.New(canvas), margin, -margin, margin/2, -margin/2))
	return canvas.Image()
}

// roundTicks places axis ticks at round values and labels them with a decimal
// comma.
type roundTicks struct{}

func (roundTicks) Ticks(low, high float64) []plot.Tick {
	step := tickStep(high - low)
	first := int(math.Ceil(low/step - tickTolerance))
	last := int(math.Floor(high/step + tickTolerance))
	ticks := make([]plot.Tick, 0, last-first+1)
	for i := first; i <= last; i++ {
		value := float64(i) * step
		ticks = append(ticks, plot.Tick{Value: value, Label: numfmt.Compact(value)})
	}
	return ticks
}

// tickStep returns the smallest round distance between ticks that puts at
// most targetTicks intervals on an axis of the given span.
func tickStep(span float64) float64 {
	raw := span / targetTicks
	magnitude := math.Pow(10, math.Floor(math.Log10(raw)))
	for _, multiplier := range tickMultipliers {
		if step := multiplier * magnitude; step >= raw {
			return step
		}
	}
	return raw
}
