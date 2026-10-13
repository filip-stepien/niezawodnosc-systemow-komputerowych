package ui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"empirical-reliability/internal/analysis"
	"empirical-reliability/internal/numfmt"
)

const sampleSeparator = "\n"

// sampleView lists the individual failure times of the analysed sample.
type sampleView struct {
	times *widget.Entry
}

func newSampleView() *sampleView {
	return &sampleView{times: multiLineEntry(placeholderSample)}
}

func (v *sampleView) tab() *container.TabItem {
	hint := wrappedLabel(fmt.Sprintf(hintSampleFormat, tabTimes))
	return container.NewTabItem(tabSample, container.NewBorder(hint, nil, nil, nil, v.times))
}

func (v *sampleView) show(a *analysis.Analysis) {
	formatted := make([]string, len(a.Result.Times))
	for i, t := range a.Result.Times {
		formatted[i] = numfmt.Compact(t)
	}
	v.times.SetText(strings.Join(formatted, sampleSeparator))
}
