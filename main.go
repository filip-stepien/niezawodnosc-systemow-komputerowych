// Command empirical-reliability computes the empirical reliability indicators
// F*(t), R*(t), f*(t), λ*(t) and E*T from a random sample or from manually
// entered data.
package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

const (
	windowTitle  = "Empiryczne wskaźniki niezawodności"
	greetingText = "Witaj, świecie!"
	windowWidth  = 480
	windowHeight = 320
)

func main() {
	window := app.New().NewWindow(windowTitle)
	window.SetContent(widget.NewLabel(greetingText))
	window.Resize(fyne.NewSize(windowWidth, windowHeight))
	window.ShowAndRun()
}
