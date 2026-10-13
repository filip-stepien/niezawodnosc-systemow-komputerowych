package ui

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

// host gives forms and views access to the window they are shown in.
type host struct {
	window fyne.Window
}

// reportError shows err to the user; a nil error is ignored.
func (h host) reportError(err error) {
	if err != nil {
		dialog.ShowError(errors.New(describe(err)), h.window)
	}
}

func (h host) inform(title, message string) {
	dialog.ShowInformation(title, message, h.window)
}
