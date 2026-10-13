package chart

import (
	"image/color"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/vg"
)

const (
	imageWidth  = 6 * vg.Inch
	imageHeight = 4 * vg.Inch
	resolution  = 150 // dots per inch

	// maxMarkers is the largest number of points still drawn with markers.
	maxMarkers    = 60
	theorySamples = 400

	// headroom leaves space above the highest value.
	headroom = 1.08
	// offChartFactor places undefined theoretical values far above the chart,
	// where they are clipped away.
	offChartFactor = 10

	// targetTicks is the largest number of intervals between ticks on an axis.
	targetTicks = 6
	// tickTolerance keeps a tick lying on the end of an axis despite
	// floating-point rounding.
	tickTolerance = 1e-9

	probabilityMax = 1.0
	// emptyRangeMax is the value range of an indicator that is zero everywhere.
	emptyRangeMax = 1.0
)

var (
	margin         = vg.Points(10)
	titlePadding   = vg.Points(6)
	empiricalWidth = vg.Points(1.6)
	theoryWidth    = vg.Points(1.4)
	markerRadius   = vg.Points(2.5)
	legendPadding  = vg.Points(4)
	legendFontSize = vg.Points(10)
	theoryDashes   = []vg.Length{vg.Points(5), vg.Points(3)}

	empiricalColor = color.RGBA{R: 0x1f, G: 0x77, B: 0xb4, A: 0xff}
	theoryColor    = color.RGBA{R: 0xd6, G: 0x27, B: 0x28, A: 0xff}

	chartFont = font.Font{Typeface: "Liberation", Variant: "Sans"}

	// tickMultipliers are the leading digits a round tick distance may have.
	tickMultipliers = []float64{1, 2, 5, 10}
)

func init() {
	// gonum/plot takes the font of every text element from this package-level default.
	plot.DefaultFont = chartFont
}
