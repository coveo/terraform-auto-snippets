package commondata

// Argument describes an argument to either provider, resource or data
type Argument struct {
	Base     `yaml:",inline"`
	Required bool
	Fields   []Argument
}
