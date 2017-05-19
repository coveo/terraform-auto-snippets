package commondata

// Provider describes a terraform provider
type Provider struct {
	Name          string       `yaml:"name"`
	Title         string       `yaml:"title"`
	Description   string       `yaml:"description"`
	URL           string       `yaml:"url"`
	Arguments     ArgumentList `yaml:"arguments"`
	DataResources ResourceList `yaml:"dataresources"`
	Resources     ResourceList `yaml:"resources"`
}

// Resource describes a terraform data or resource
type Resource struct {
	Name        string       `yaml:"name"`
	Description string       `yaml:"description"`
	Section     string       `yaml:"section"`
	URL         string       `yaml:"url"`
	Arguments   ArgumentList `yaml:"arguments"`
}

// Argument describes an argument to either provider, resource or data
type Argument struct {
	Name        string       `yaml:"name"`
	Description string       `yaml:"description"`
	URL         string       `yaml:"url"`
	Required    bool         `yaml:"required"`
	Fields      ArgumentList `yaml:"fields"`
}

// ProviderMap holds a map of all providers
type ProviderMap map[string]Provider

// ResourceList hosts a list of resources or data for a provider
type ResourceList []Resource

// ArgumentList hosts a list of arguments for a provider, a resource or a data
type ArgumentList []Argument
