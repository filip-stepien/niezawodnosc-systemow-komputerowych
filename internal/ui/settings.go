package ui

const (
	appID = "empirical-reliability"

	windowWidth  float32 = 1280
	windowHeight float32 = 780
	// inputPaneShare is the part of the window width taken by the input forms.
	inputPaneShare = 0.34

	// secondaryActionColumns is the number of secondary buttons shown side by side.
	secondaryActionColumns = 2

	chartColumns           = 2
	chartMinWidth  float32 = 300
	chartMinHeight float32 = 200

	// Decimal places of the displayed values.
	probabilityDecimals = 4
	rateDecimals        = 5
	meanDecimals        = 4
	exportMeanDecimals  = 6

	defaultRandomSize = "1000"

	resultsFileName    = "results.csv"
	chartFileExtension = ".png"
)
