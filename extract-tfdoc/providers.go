package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"net/url"
	"sort"
	"strings"
	"sync"

	data "github.com/coveo/terraform-auto-snippets/common_data"
	"github.com/coveo/terraform-auto-snippets/utils"
)

const terraformDocURL = "https://www.terraform.io/docs/providers/index.html"

// ParseTerraformDocumentation launch the process of parsing the terraform documentation web site
func ParseTerraformDocumentation(nbWorkers int) ProviderWorkForce {
	return ProviderWorkForce{utils.StartDocumentWorkers(nbWorkers), nil}
}

// ProviderWorkForce is an implementation of a work force
type ProviderWorkForce struct {
	*utils.WorkForce
	rootURL *url.URL
}

// GetProviders returns the list of all providers matching the provided filters
func (wf ProviderWorkForce) GetProviders(filters ...string) (providers data.ProviderList, err error) {
	wf.rootURL = wf.resolve(nil, terraformDocURL)
	utils.Assert(wf.rootURL != nil, "Invalid URI %s", terraformDocURL)
	filters = utils.Expand(",; ", filters)

	var (
		mutex          sync.Mutex
		subTasks       sync.WaitGroup
		totalData      int
		totalResources int
	)

	err = wf.ProcessDocument(*wf.rootURL, func(doc *goquery.Document, responseChannel chan error) {
		var err error
		defer func() { responseChannel <- err }()

		doc.Find(".active ul a").Each(func(i int, s *goquery.Selection) {
			name := s.Text()
			if href, ok := s.Attr("href"); ok {
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

				subTasks.Add(1)
				go func() {
					defer subTasks.Done()
					if provider, err := wf.getProvider(href); err == nil {
						defer mutex.Unlock()
						mutex.Lock()
						providers = append(providers, provider)
						totalData += len(provider.DataResources)
						totalResources += len(provider.Resources)

						if len(provider.DataResources)+len(provider.Resources) == 0 {
							utils.PrintError("No resource found for %s", name)
						}
					} else {
						utils.PrintError("Unable to get provider %s: %v", name, err)
					}
				}()
			}
		})

		return
	})

	subTasks.Wait()
	utils.PrintMessage("%d providers, %d resources, %d data sources", len(providers), totalResources, totalData)
	sort.Sort(data.ProvidersCompleteness{ProviderList: providers})

	return
}

// getProvider returns details about the requested provider
func (wf ProviderWorkForce) getProvider(link string) (provider *data.Provider, err error) {
	uri := wf.resolve(wf.rootURL, link)
	if uri == nil {
		err = fmt.Errorf("Malformed URI %s", link)
		return
	}

	err = wf.ProcessDocument(*uri, func(doc *goquery.Document, responseChannel chan error) {
		var err error
		defer func() { responseChannel <- err }()

		title := doc.Find("h1")
		switch title.Length() {
		case 0:
			err = fmt.Errorf("No title found in %s", uri.String())
			return
		case 1:
			break
		default:
			utils.PrintWarning("Found more that one title (%d) in %s", title.Length(), uri.String())
			title = title.First()
		}

		titleNode := title.Nodes[0].LastChild
		if titleNode.Type != html.TextNode {
			err = fmt.Errorf("Malformed title node in %s %v", uri.String(), title.Text())
			return
		}

		providerName := utils.Trim(titleNode.Data)

		id, ok := title.Attr("id")
		if !ok {
			utils.PrintWarning("No id found in title for %s", uri.String())
			id = strings.ToLower(strings.Replace(providerName, " ", "-", -1))
		}

		resources := resourceWorkForce{wf, providerName, *uri, doc.Find("ul.nav.docs-sidenav li.active")}

		provider = &data.Provider{
			Base: data.Base{
				Name:        strings.Replace(id, "-provider", "", 1),
				Description: title.Next().Text(),
				URL:         uri.String(),
			},
			Title:         providerName,
			Arguments:     getArgs(providerName, *uri, doc.Find("#argument-reference")),
			DataResources: resources.get(true),
			Resources:     resources.get(false),
		}

		doc.Find("pre.highlight.hcl").Each(func(i int, example *goquery.Selection) {
			provider.Examples = append(provider.Examples, example.Text())
		})

		utils.PrintInfo(" ▪︎ %-33s %s", providerName, utils.MessagePrinter(" - %3d resource(s), %2d data", len(provider.Resources), len(provider.DataResources)))
		return
	})

	return
}

func (wf ProviderWorkForce) resolve(uri *url.URL, link string) *url.URL {
	linkURI, err := url.Parse(link)
	if err != nil {
		return nil
	}

	if uri == nil {
		return linkURI
	}
	return uri.ResolveReference(linkURI)
}
