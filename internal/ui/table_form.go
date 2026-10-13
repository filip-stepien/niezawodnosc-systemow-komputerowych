package ui

import (
	"strconv"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"empirical-reliability/internal/analysis"
	"empirical-reliability/internal/input"
	"empirical-reliability/internal/numfmt"
)

// tableForm collects grouped data: the sample size and a failure table.
type tableForm struct {
	size   *widget.Entry
	data   *widget.Entry
	submit *widget.Button
}

func newTableForm(run runner) *tableForm {
	f := &tableForm{
		size: widget.NewEntry(),
		data: widget.NewMultiLineEntry(),
	}
	f.submit = primaryButton(buttonCompute, func() { run(f.analyze) })
	f.loadExample()
	return f
}

func (f *tableForm) tab() *container.TabItem {
	top := container.NewVBox(formGrid(labelled(labelSampleSize, f.size)), wrappedLabel(hintTable))
	bottom := container.NewVBox(
		f.submit,
		container.NewGridWithColumns(secondaryActionColumns,
			widget.NewButton(buttonExample, f.loadExample),
			widget.NewButton(buttonClear, f.clear),
		),
	)
	return container.NewTabItem(tabTable, container.NewBorder(top, bottom, nil, nil, f.data))
}

func (f *tableForm) loadExample() {
	f.size.SetText(strconv.Itoa(exampleSampleSize))
	f.data.SetText(input.FormatTable(exampleObservations))
}

func (f *tableForm) clear() {
	f.size.SetText("")
	f.data.SetText("")
}

func (f *tableForm) analyze() (*analysis.Analysis, error) {
	sampleSize, err := parseField(labelSampleSize, f.size.Text, numfmt.ParseInt)
	if err != nil {
		return nil, err
	}
	observations, err := input.ParseTable(f.data.Text)
	if err != nil {
		return nil, err
	}
	return analysis.FromTable(sampleSize, observations)
}
