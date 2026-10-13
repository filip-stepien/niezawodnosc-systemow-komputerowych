package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// lightTheme forces the light variant of the default theme, so that the charts
// drawn on a white background match the rest of the window regardless of the
// system settings.
type lightTheme struct {
	fyne.Theme
}

func newLightTheme() fyne.Theme {
	return lightTheme{Theme: theme.DefaultTheme()}
}

func (t lightTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	return t.Theme.Color(name, theme.VariantLight)
}
