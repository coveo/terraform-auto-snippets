package commondata

// Resource describes a terraform data or resource
type Resource struct {
	Base      `yaml:",inline"`
	Section   string
	Arguments []Argument
	Examples  []string
}
