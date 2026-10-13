package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// tabbed is a part of the window shown as a tab.
type tabbed interface {
	tab() *container.TabItem
}

func tabsOf(parts ...tabbed) *container.AppTabs {
	tabs := container.NewAppTabs()
	for _, part := range parts {
		tabs.Append(part.tab())
	}
	return tabs
}

// newSelect builds a select widget with an explicit placeholder.
// widget.NewSelect is avoided because it looks up a default placeholder that
// Fyne has no Polish translation for, and logs an error on every call.
func newSelect(options []string, placeholder string, changed func(string)) *widget.Select {
	selector := &widget.Select{Options: options, PlaceHolder: placeholder, OnChanged: changed}
	selector.ExtendBaseWidget(selector)
	return selector
}

func entryWithText(text string) *widget.Entry {
	entry := widget.NewEntry()
	entry.SetText(text)
	return entry
}

func entryWithPlaceholder(placeholder string) *widget.Entry {
	entry := widget.NewEntry()
	entry.SetPlaceHolder(placeholder)
	return entry
}

func multiLineEntry(placeholder string) *widget.Entry {
	entry := widget.NewMultiLineEntry()
	entry.SetPlaceHolder(placeholder)
	entry.Wrapping = fyne.TextWrapWord
	return entry
}

func wrappedLabel(text string) *widget.Label {
	label := widget.NewLabel(text)
	label.Wrapping = fyne.TextWrapWord
	return label
}

func primaryButton(label string, action func()) *widget.Button {
	button := widget.NewButton(label, action)
	button.Importance = widget.HighImportance
	return button
}

// labelled returns a label followed by the field it describes, as one row of
// a form grid.
func labelled(label string, field fyne.CanvasObject) []fyne.CanvasObject {
	return []fyne.CanvasObject{widget.NewLabel(label), field}
}

func formGrid(rows ...[]fyne.CanvasObject) *fyne.Container {
	grid := container.New(layout.NewFormLayout())
	for _, row := range rows {
		grid.Objects = append(grid.Objects, row...)
	}
	return grid
}
