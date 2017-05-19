package common_data

type Provider struct {
	Name          string     `yaml:"name"`
	Title         string     `yaml:"title"`
	Description   string     `yaml:"description"`
	URL           string     `yaml:"url"`
	Arguments     []Argument `yaml:"arguments"`
	DataResources []Resource `yaml:"dataresources"`
	Resources     []Resource `yaml:"resources"`
}

type Argument struct {
	Name        string     `yaml:"name"`
	Description string     `yaml:"description"`
	URL         string     `yaml:"url"`
	Required    bool       `yaml:"required"`
	Arguments   []Argument `yaml:"arguments"`
}

type Resource struct {
	Name        string     `yaml:"name"`
	Description string     `yaml:"description"`
	URL         string     `yaml:"url"`
	Arguments   []Argument `yaml:"arguments"`
}
