package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"net/url"
	"strings"

	data "github.com/coveo/terraform-auto-snippets/common_data"
)

func getProviders(uri url.URL, filters ...string) (providers data.ProviderMap, err error) {
	providers = map[string]data.Provider{}
	doc, err := getDocument(uri)
	if err != nil {
		return
	}

	var totalData, totalResources int
	doc.Find(".active ul a").Each(func(i int, s *goquery.Selection) {
		name := s.Text()
		href, ok := s.Attr("href")
		if ok {
			link, err := url.Parse(href)
			if err != nil {
				PrintWarning("Malformed URL %v for provider %s", href, name)
				return
			}

			// We consider that there is a match if there is no filter
			match := len(filters) == 0
			for _, filter := range filters {
				match = strings.Contains(strings.ToLower(name), strings.ToLower(filter))
				if match {
					break
				}
			}

			if !match {
				// No match, we skip the provider
				return
			}

			if provider, err := getProvider(*uri.ResolveReference(link)); err == nil {
				providers[name] = *provider
				totalData += len(provider.DataResources)
				totalResources += len(provider.Resources)

				if len(provider.DataResources)+len(provider.Resources) == 0 {
					PrintError("No resource found for %s", name)
				}
			} else {
				PrintError("Unable to get provider %s: %v", name, err)
			}
		}
	})

	PrintInfo("%d providers, %d data sources, %d resources", len(providers), totalData, totalResources)
	return
}

func getProvider(uri url.URL) (provider *data.Provider, err error) {
	doc, err := getDocument(uri)
	if err != nil {
		return
	}

	title := doc.Find("h1")
	switch title.Length() {
	case 0:
		err = fmt.Errorf("No title found in %s", uri.String())
		return
	case 1:
		break
	default:
		PrintWarning("Found more that one title (%d) in %s", title.Length(), uri.String())
		title = title.First()
	}

	titleNode := title.Nodes[0].LastChild
	if titleNode.Type != html.TextNode {
		err = fmt.Errorf("Malformed title node in %s %v", uri.String(), title.Text())
		return
	}

	providerName := trim(titleNode.Data)

	id, ok := title.Attr("id")
	if !ok {
		PrintWarning("No id found in title for %s", uri.String())
		id = strings.ToLower(strings.Replace(providerName, " ", "-", -1))
	}

	PrintInfo(providerName)
	provider = &data.Provider{
		Name:          strings.Replace(id, "-provider", "", 1),
		Title:         providerName,
		Description:   title.Next().Text(),
		URL:           uri.String(),
		Arguments:     getArgs(providerName, uri, doc.Find("#argument-reference")),
		DataResources: getResources(providerName, uri, doc.Find("ul.nav.docs-sidenav li.active"), true),
		Resources:     getResources(providerName, uri, doc.Find("ul.nav.docs-sidenav li.active"), false),
	}
	fmt.Println()

	return
}
