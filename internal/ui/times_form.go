package ui

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"empirical-reliability/internal/analysis"
	"empirical-reliability/internal/input"
	"empirical-reliability/internal/numfmt"
)

// timesForm collects individual failure times.
type timesForm struct {
	size   *widget.Entry
	step   *widget.Entry
	data   *widget.Entry
	submit *widget.Button
}

func newTimesForm(run runner) *timesForm {
	f := &timesForm{
		size: entryWithPlaceholder(placeholderTimesCount),
		step: entryWithPlaceholder(placeholderAutomatic),
		data: multiLineEntry(placeholderTimes),
	}
	f.submit = primaryButton(buttonCompute, func() { run(f.analyze) })
	return f
}

func (f *timesForm) tab() *container.TabItem {
	top := container.NewVBox(
		formGrid(labelled(labelSampleSize, f.size), labelled(labelStep, f.step)),
		wrappedLabel(hintTimes),
	)
	bottom := container.NewVBox(f.submit, widget.NewButton(buttonClear, f.clear))
	return container.NewTabItem(tabTimes, container.NewBorder(top, bottom, nil, nil, f.data))
}

func (f *timesForm) clear() {
	f.data.SetText("")
}

func (f *timesForm) analyze() (*analysis.Analysis, error) {
	times, err := input.ParseTimes(f.data.Text)
	if err != nil {
		return nil, err
	}
	sampleSize, err := f.sampleSize(len(times))
	if err != nil {
		return nil, err
	}
	step, err := parseStep(f.step.Text)
	if err != nil {
		return nil, err
	}
	return analysis.FromTimes(times, sampleSize, step)
}

// sampleSize returns the entered sample size, or the number of failure times
// when the field is left blank.
func (f *timesForm) sampleSize(timeCount int) (int, error) {
	if isBlank(f.size.Text) {
		return timeCount, nil
	}
	return parseField(labelSampleSize, f.size.Text, numfmt.ParseInt)
}
