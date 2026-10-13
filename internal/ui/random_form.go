package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"empirical-reliability/internal/analysis"
	"empirical-reliability/internal/distribution"
	"empirical-reliability/internal/numfmt"
)

// randomForm collects the settings of a random sample to draw.
type randomForm struct {
	specs []distribution.Spec

	kind       *widget.Select
	parameters []*widget.Entry // one entry per parameter of the selected distribution
	size       *widget.Entry
	step       *widget.Entry
	seed       *widget.Entry
	fields     *fyne.Container
	submit     *widget.Button
	content    *fyne.Container
}

func newRandomForm(specs []distribution.Spec, run runner) *randomForm {
	f := &randomForm{
		specs:  specs,
		size:   entryWithText(defaultRandomSize),
		step:   entryWithPlaceholder(placeholderAutomatic),
		seed:   entryWithPlaceholder(placeholderRandomSeed),
		fields: formGrid(),
	}
	f.kind = newSelect(f.distributionNames(), placeholderDistribution, func(string) { f.showParameters() })
	f.submit = primaryButton(buttonDraw, func() { run(f.analyze) })
	f.content = container.NewVBox(f.fields, f.submit, wrappedLabel(hintRandom))
	f.kind.SetSelectedIndex(0)
	return f
}

func (f *randomForm) tab() *container.TabItem {
	return container.NewTabItem(tabRandom, f.content)
}

func (f *randomForm) distributionNames() []string {
	names := make([]string, len(f.specs))
	for i, spec := range f.specs {
		names[i] = distributionNames[spec.ID]
	}
	return names
}

func (f *randomForm) selected() distribution.Spec {
	return f.specs[f.kind.SelectedIndex()]
}

// showParameters rebuilds the form so that it holds one field per parameter
// of the selected distribution, filled with the default value.
func (f *randomForm) showParameters() {
	spec := f.selected()
	f.parameters = make([]*widget.Entry, len(spec.Parameters))
	rows := [][]fyne.CanvasObject{labelled(labelDistribution, f.kind)}
	for i, parameter := range spec.Parameters {
		f.parameters[i] = entryWithText(numfmt.Compact(parameter.Default))
		rows = append(rows, labelled(parameterLabel(parameter.Symbol), f.parameters[i]))
	}
	rows = append(rows,
		labelled(labelRandomSize, f.size),
		labelled(labelStep, f.step),
		labelled(labelSeed, f.seed),
	)
	f.fields.Objects = formGrid(rows...).Objects
	// The number of rows changes, so everything below the fields moves.
	f.content.Refresh()
}

func (f *randomForm) analyze() (*analysis.Analysis, error) {
	spec := f.selected()
	parameters, err := f.parameterValues(spec)
	if err != nil {
		return nil, err
	}
	size, err := parseField(labelRandomSize, f.size.Text, numfmt.ParseInt)
	if err != nil {
		return nil, err
	}
	step, err := parseStep(f.step.Text)
	if err != nil {
		return nil, err
	}
	seed, err := parseSeed(f.seed.Text)
	if err != nil {
		return nil, err
	}
	return analysis.FromRandom(analysis.RandomRequest{
		Spec:       spec,
		Parameters: parameters,
		Size:       size,
		Seed:       seed,
		Step:       step,
	})
}

func (f *randomForm) parameterValues(spec distribution.Spec) ([]float64, error) {
	values := make([]float64, len(f.parameters))
	for i, entry := range f.parameters {
		value, err := parseField(parameterLabel(spec.Parameters[i].Symbol), entry.Text, numfmt.ParseFloat)
		if err != nil {
			return nil, err
		}
		values[i] = value
	}
	return values, nil
}

func parameterLabel(symbol string) string {
	return fmt.Sprintf(parameterLabelFormat, symbol, parameterNames[symbol])
}
