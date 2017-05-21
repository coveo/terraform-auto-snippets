package commondata

// Provider describes a terraform provider
type Provider struct {
	Base          `yaml:",inline"`
	Title         string
	Arguments     []Argument
	DataResources []Resource
	Resources     []Resource
}

// Count returns the total number of resources (data + resource) in the provider
func (p Provider) Count() int { return len(p.DataResources) + len(p.Resources) }

// MarshalYAML implements custom marshaler for Provider in order to serialize the count property along with the object
func (p Provider) MarshalYAML() (interface{}, error) {
	type Base Provider // This is required to avoid recursion. Base does not implement MarshalYAML while Provider does.
	type copy struct {
		Base  `yaml:",inline"`
		Count int
	}
	return copy{Base(p), p.Count()}, nil
}

// ProviderList holds a list of all providers
type ProviderList []*Provider

// Implement sort on provider list
func (p ProviderList) Len() int           { return len(p) }
func (p ProviderList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p ProviderList) Less(i, j int) bool { return p[i].Name < p[j].Name }

// ProvidersCompleteness allows sorting by completeness (i.e. providers which have more resources come first)
type ProvidersCompleteness struct{ ProviderList }

// Less defines the implementation of the sort comparison operator
func (p ProvidersCompleteness) Less(i, j int) bool {
	pi, pj := p.ProviderList[i], p.ProviderList[j]
	return pi.Count() > pj.Count() || pi.Count() == pj.Count() && pi.Name < pj.Name
}
