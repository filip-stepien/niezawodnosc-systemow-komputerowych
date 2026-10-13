package distribution

// Parameter describes one parameter of a distribution family.
type Parameter struct {
	Symbol  string
	Default float64
}

// Spec describes a distribution family: its identifier, its parameters and
// how to build a distribution from parameter values. Supporting a new family
// only takes a new Spec in the catalog.
type Spec struct {
	ID         string
	Parameters []Parameter
	build      func(values []float64) Distribution
}

// New builds a distribution from one value per parameter, in the order of
// Parameters. Every parameter must be positive.
func (s Spec) New(values ...float64) (Distribution, error) {
	if len(values) != len(s.Parameters) {
		return nil, ErrParameterCount
	}
	for i, value := range values {
		if !isPositive(value) {
			return nil, &ParameterError{Symbol: s.Parameters[i].Symbol}
		}
	}
	return s.build(values), nil
}

// The distribution families available to the application.
var (
	Exponential = Spec{
		ID:         "exponential",
		Parameters: []Parameter{{Symbol: "λ", Default: 0.01}},
		build:      func(v []float64) Distribution { return exponential{rate: v[0]} },
	}
	Weibull = Spec{
		ID:         "weibull",
		Parameters: []Parameter{{Symbol: "β", Default: 2}, {Symbol: "η", Default: 100}},
		build:      func(v []float64) Distribution { return weibull{shape: v[0], scale: v[1]} },
	}
	Normal = Spec{
		ID:         "normal",
		Parameters: []Parameter{{Symbol: "μ", Default: 100}, {Symbol: "σ", Default: 20}},
		build:      func(v []float64) Distribution { return newNormal(v[0], v[1]) },
	}
)

// Catalog lists the available distribution families.
func Catalog() []Spec {
	return []Spec{Exponential, Weibull, Normal}
}
