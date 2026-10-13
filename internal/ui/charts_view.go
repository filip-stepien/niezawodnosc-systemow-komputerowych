package ui

import (
	"fmt"
	"image"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"empirical-reliability/internal/analysis"
	"empirical-reliability/internal/chart"
	"empirical-reliability/internal/export"
)

// chartsView shows one chart per reliability indicator.
type chartsView struct {
	host     host
	panels   map[chart.Indicator]*canvas.Image
	rendered map[chart.Indicator]image.Image // empty until something is computed
}

func newChartsView(h host) *chartsView {
	v := &chartsView{host: h, panels: map[chart.Indicator]*canvas.Image{}}
	for _, indicator := range chart.Indicators() {
		v.panels[indicator] = newChartPanel()
	}
	return v
}

func newChartPanel() *canvas.Image {
	panel := canvas.NewImageFromImage(nil)
	panel.FillMode = canvas.ImageFillContain
	panel.SetMinSize(fyne.NewSize(chartMinWidth, chartMinHeight))
	return panel
}

func (v *chartsView) tab() *container.TabItem {
	grid := container.NewGridWithColumns(chartColumns)
	for _, indicator := range chart.Indicators() {
		grid.Add(v.panels[indicator])
	}
	save := widget.NewButton(buttonSaveCharts, v.saveCharts)
	return container.NewTabItem(tabCharts, container.NewBorder(nil, save, nil, nil, grid))
}

func (v *chartsView) show(a *analysis.Analysis) {
	rendered, err := renderCharts(a)
	if err != nil {
		v.host.reportError(err)
		return
	}
	v.rendered = rendered
	for indicator, img := range rendered {
		v.panels[indicator].Image = img
		v.panels[indicator].Refresh()
	}
}

func renderCharts(a *analysis.Analysis) (map[chart.Indicator]image.Image, error) {
	rendered := map[chart.Indicator]image.Image{}
	for _, indicator := range chart.Indicators() {
		img, err := chart.Render(indicator, a.Result, a.Model(), chartLabels(indicator))
		if err != nil {
			return nil, err
		}
		rendered[indicator] = img
	}
	return rendered, nil
}

func (v *chartsView) saveCharts() {
	if len(v.rendered) == 0 {
		v.host.reportError(errNothingToSave)
		return
	}
	dialog.ShowFolderOpen(func(folder fyne.ListableURI, err error) {
		if err != nil {
			v.host.reportError(err)
			return
		}
		if folder == nil { // the dialog was cancelled
			return
		}
		if err := v.saveTo(folder.Path()); err != nil {
			v.host.reportError(err)
			return
		}
		v.host.inform(dialogSavedTitle, fmt.Sprintf(dialogChartsSavedFormat, folder.Path()))
	}, v.host.window)
}

func (v *chartsView) saveTo(folder string) error {
	for indicator, img := range v.rendered {
		path := filepath.Join(folder, indicator.Slug()+chartFileExtension)
		if err := export.SavePNG(path, img); err != nil {
			return err
		}
	}
	return nil
}
