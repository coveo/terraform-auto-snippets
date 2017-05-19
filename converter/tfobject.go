package main

type TfObject struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

type TfObjectUrl struct {
	TfObject `yaml:",inline"`
	Url      string `yaml:"url"`
}

type TfObjectResource struct {
	TfObjectUrl `yaml:",inline"`
	Arguments   []Argument `yaml:"arguments"`
}

type Provider struct {
	TfObjectResource `yaml:",inline"`
	Dataresources    []Dataresource `yaml:"dataresources"`
	Resources        []Resource     `yaml:"resources"`
}

type Argument struct {
	TfObject `yaml:",inline"`
	Requierd string `yaml:"required"`
}

type Resource struct {
	TfObjectResource `yaml:",inline"`
}

type Dataresource struct {
	TfObjectResource `yaml:",inline"`
}
