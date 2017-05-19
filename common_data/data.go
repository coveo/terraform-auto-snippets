package common_data

type Provider struct {
	Name          string     `yaml:"name"`
	Title         string     `yaml:"title"`
	Description   string     `yaml:"description"`
	URL           string     `yaml:"url"`
	Arguments     []Argument `yaml:"arguments"`
	DataResources []Data     `yaml:"dataresources"`
	Resources     []Resource `yaml:"resources"`
}

type Argument struct {
	Name        string     `yaml:"name"`
	Description string     `yaml:"description"`
	Required    bool       `yaml:"required"`
	Arguments   []Argument `yaml:"arguments"`
}

type Data struct {
	Name        string     `yaml:"name"`
	Description string     `yaml:"description"`
	URL         string     `yaml:"url"`
	Arguments   []Argument `yaml:"arguments"`
}

type Resource Data
