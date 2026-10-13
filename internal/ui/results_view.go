package ui

import (
	"errors"
	"io"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"empirical-reliability/internal/analysis"
	"empirical-reliability/internal/export"
	"empirical-reliability/internal/reliability"
)

// resultsView shows the summary and the table of indicators.
type resultsView struct {
	host    host
	current *analysis.Analysis // nil until something is computed

	summary *widget.Label
	table   *widget.Table
	content *fyne.Container
}

func newResultsView(h host) *resultsView {
	v := &resultsView{host: h, summary: wrappedLabel(summaryPrompt)}
	v.table = v.newTable()
	save := widget.NewButton(buttonSaveCSV, v.saveCSV)
	v.content = container.NewBorder(v.summary, save, nil, nil, v.table)
	return v
}

func (v *resultsView) tab() *container.TabItem {
	return container.NewTabItem(tabResults, v.content)
}

func (v *resultsView) show(a *analysis.Analysis) {
	v.current = a
	v.summary.SetText(summaryText(a))
	v.table.Refresh()
	v.table.ScrollToTop()
	// The summary changes its height, so the table has to be laid out again.
	v.content.Refresh()
}

func (v *resultsView) newTable() *widget.Table {
	table := widget.NewTable(
		func() (int, int) { return len(v.rows()), len(resultColumns) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			row := v.rows()[id.Row]
			cell.(*widget.Label).SetText(resultColumns[id.Col].cell(id.Row, row))
		},
	)
	table.ShowHeaderRow = true
	table.CreateHeader = newHeaderCell
	table.UpdateHeader = func(id widget.TableCellID, cell fyne.CanvasObject) {
		cell.(*widget.Label).SetText(resultColumns[id.Col].header)
	}
	for i, column := range resultColumns {
		table.SetColumnWidth(i, column.width)
	}
	return table
}

func newHeaderCell() fyne.CanvasObject {
	label := widget.NewLabel("")
	label.TextStyle.Bold = true
	return label
}

func (v *resultsView) rows() []reliability.Row {
	if v.current == nil {
		return nil
	}
	return v.current.Result.Rows
}

func (v *resultsView) saveCSV() {
	if v.current == nil {
		v.host.reportError(errNothingToSave)
		return
	}
	save := dialog.NewFileSave(func(target fyne.URIWriteCloser, err error) {
		if err != nil {
			v.host.reportError(err)
			return
		}
		if target == nil { // the dialog was cancelled
			return
		}
		v.host.reportError(writeResults(target, v.current.Result))
	}, v.host.window)
	save.SetFileName(resultsFileName)
	readableFileDialog(save, v.host.window).Show()
}

func writeResults(target io.WriteCloser, result *reliability.Result) error {
	return errors.Join(export.WriteCSV(target, exportRecords(result)), target.Close())
}
