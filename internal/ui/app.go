// Package ui is the graphical user interface of the application.
package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

// Run opens the main window and blocks until it is closed.
func Run() {
	application := app.NewWithID(appID)
	application.Settings().SetTheme(newLightTheme())

	window := application.NewWindow(windowTitle)
	window.SetContent(newMainWindow(window).content())
	window.Resize(fyne.NewSize(windowWidth, windowHeight))
	window.ShowAndRun()
}
