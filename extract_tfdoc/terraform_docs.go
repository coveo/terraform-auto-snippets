package main

import (
	"net/url"

	data "github.com/coveo/terraform-auto-snippets/common_data"
	"github.com/coveo/terraform-auto-snippets/utils"
)

const terraformDocURL = "https://www.terraform.io/docs/providers/index.html"

func getTerraformDocs(filters ...string) (providers data.ProviderMap, err error) {
	uri, err := url.Parse(terraformDocURL)
	if err == nil {
		providers, err = getProviders(*uri, utils.Expand(",; ", filters))
	}
	return
}
