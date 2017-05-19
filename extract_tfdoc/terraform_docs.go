package main

import (
	data "github.com/coveo/terraform-auto-snippets/common_data"
	"net/url"
)

const terraformDocURL = "https://www.terraform.io/docs/providers/index.html"

func getTerraformDocs() (providers data.ProviderMap, err error) {
	uri, err := url.Parse(terraformDocURL)
	if err == nil {
		providers, err = getProviders(*uri)
	}
	return
}
