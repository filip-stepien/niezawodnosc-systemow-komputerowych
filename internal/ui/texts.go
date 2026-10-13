package ui

import (
	"empirical-reliability/internal/chart"
	"empirical-reliability/internal/distribution"
)

// All texts shown to the user are collected in this file. The application is
// presented in Polish; the rest of the code base is in English.

const (
	windowTitle = "Empiryczne wskaźniki niezawodności"

	tabRandom  = "Próba losowa"
	tabTable   = "Tabela uszkodzeń"
	tabTimes   = "Lista czasów"
	tabResults = "Wyniki"
	tabCharts  = "Wykresy"
	tabSample  = "Próba"

	labelDistribution    = "Rozkład trwałości"
	labelRandomSize      = "Liczność próby n"
	labelSampleSize      = "Liczność próby n(0)"
	labelStep            = "Długość przedziału Δt"
	labelSeed            = "Ziarno generatora"
	parameterLabelFormat = "%s (%s)"

	placeholderDistribution = "wybierz rozkład"
	placeholderAutomatic    = "automatycznie"
	placeholderRandomSeed   = "losowe"
	placeholderTimesCount   = "liczba podanych czasów"
	placeholderTimes        = "np. 12,5 30 41,2 55 ..."
	placeholderSample       = "Tu pojawią się pojedyncze czasy uszkodzeń (posortowane rosnąco)."

	buttonDraw       = "Losuj i oblicz"
	buttonCompute    = "Oblicz"
	buttonExample    = "Przykład z książki"
	buttonClear      = "Wyczyść"
	buttonSaveCSV    = "Zapisz wyniki do pliku CSV"
	buttonSaveCharts = "Zapisz wykresy jako pliki PNG"

	hintRandom = "Program losuje n czasów do uszkodzenia z wybranego rozkładu, " +
		"grupuje je w przedziały o długości Δt i wyznacza wskaźniki empiryczne. " +
		"To samo ziarno daje zawsze tę samą próbę."
	hintTable = "W każdym wierszu: chwila t oraz liczba elementów uszkodzonych m(t) " +
		"zarejestrowana w tej chwili, np. „2 3”."
	hintTimes = "Czasy do uszkodzenia kolejnych elementów, rozdzielone spacjami, " +
		"średnikami lub znakami nowej linii. Jeżeli n(0) jest większe od liczby czasów, " +
		"pozostałe elementy są traktowane jako zdatne do końca obserwacji."
	hintSampleFormat = "Czasy uszkodzeń z próby, posortowane rosnąco. " +
		"Można je skopiować do zakładki „%s”."

	summaryPrompt                    = "Wprowadź dane po lewej stronie i wybierz „" + buttonCompute + "”."
	summarySampleFormat              = "Liczność próby n(0) = %d, liczba przedziałów: %d, czas obserwacji: %s"
	summaryCountsFormat              = "Elementy uszkodzone: %d, zdatne po zakończeniu obserwacji: %d"
	summaryMeanFormat                = "E*T ze wzoru (5-15) = %s"
	summaryMeanFromInspectionsFormat = "E*T ze wzoru (5-15) = %s (za trwałość elementu przyjęto chwilę rejestracji uszkodzenia)"
	summaryMeanUnknown               = "E*T ze wzoru (5-15): nie można wyznaczyć, ponieważ nie wszystkie elementy uległy uszkodzeniu"
	summaryMeanFromReliabilityFormat = "E*T z funkcji R*(t) = %s"
	summaryMeanLowerBoundFormat      = "E*T z funkcji R*(t) = %s (oszacowanie od dołu)"
	summaryDistributionFormat        = "Rozkład %s, %s; ziarno generatora: %d"
	summaryParameterFormat           = "%s = %s"
	summaryParameterSeparator        = ", "
	summaryTheoreticalMeanFormat     = "Teoretyczna wartość oczekiwana ET = %s"

	exportSampleSize          = "n(0)"
	exportMean                = "E*T (5-15)"
	exportMeanFromReliability = "E*T z R*(t)"

	dialogSavedTitle        = "Zapisano"
	dialogChartsSavedFormat = "Wykresy zapisano w folderze:\n%s"

	chartTimeAxis    = "t"
	chartEmpirical   = "empiryczna"
	chartTheoretical = "teoretyczna"
)

var (
	distributionNames = map[string]string{
		distribution.Exponential.ID: "wykładniczy",
		distribution.Weibull.ID:     "Weibulla",
		distribution.Normal.ID:      "normalny",
	}

	// parameterNames describes distribution parameters by their symbol.
	parameterNames = map[string]string{
		"λ": "intensywność uszkodzeń",
		"β": "kształt",
		"η": "skala",
		"μ": "wartość średnia",
		"σ": "odchylenie standardowe",
	}

	chartTitles = map[chart.Indicator]string{
		chart.CDF:         "Dystrybuanta empiryczna F*(t)",
		chart.Reliability: "Funkcja niezawodności R*(t)",
		chart.Density:     "Gęstość prawdopodobieństwa f*(t)",
		chart.Hazard:      "Intensywność uszkodzeń λ*(t)",
	}
)

func chartLabels(indicator chart.Indicator) chart.Labels {
	return chart.Labels{
		Title:       chartTitles[indicator],
		TimeAxis:    chartTimeAxis,
		Empirical:   chartEmpirical,
		Theoretical: chartTheoretical,
	}
}
