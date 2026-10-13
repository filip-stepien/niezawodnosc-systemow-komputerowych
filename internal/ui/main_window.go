package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"

	"empirical-reliability/internal/analysis"
	"empirical-reliability/internal/distribution"
)

// view presents one aspect of an analysis.
type view interface {
	tabbed
	show(a *analysis.Analysis)
}

// mainWindow wires the input forms to the views that present their results.
type mainWindow struct {
	host host

	random *randomForm
	table  *tableForm
	times  *timesForm

	results *resultsView
	charts  *chartsView
	sample  *sampleView

	input  *container.AppTabs
	output *container.AppTabs
}

func newMainWindow(window fyne.Window) *mainWindow {
	m := &mainWindow{host: host{window: window}}

	m.random = newRandomForm(distribution.Catalog(), m.run)
	m.table = newTableForm(m.run)
	m.times = newTimesForm(m.run)
	m.input = tabsOf(m.random, m.table, m.times)

	m.results = newResultsView(m.host)
	m.charts = newChartsView(m.host)
	m.sample = newSampleView()
	m.output = tabsOf(m.results, m.charts, m.sample)

	return m
}

func (m *mainWindow) content() fyne.CanvasObject {
	split := container.NewHSplit(m.input, m.output)
	split.Offset = inputPaneShare
	return split
}

func (m *mainWindow) views() []view {
	return []view{m.results, m.charts, m.sample}
}

// run executes a computation requested by a form and presents its outcome.
func (m *mainWindow) run(c compute) {
	a, err := c()
	if err != nil {
		m.host.reportError(err)
		return
	}
	for _, v := range m.views() {
		v.show(a)
	}
}
