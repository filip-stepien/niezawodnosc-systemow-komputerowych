package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// lightTheme forces the light variant of the default theme. It is only used
// in tests, so that golden screenshots stay the same regardless of the
// system theme the tests happen to run under.
type lightTheme struct {
	fyne.Theme
}

func newLightTheme() fyne.Theme {
	return lightTheme{Theme: theme.DefaultTheme()}
}

func (t lightTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	return t.Theme.Color(name, theme.VariantLight)
}
