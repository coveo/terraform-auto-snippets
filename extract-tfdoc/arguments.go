package main

import (
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"strings"

	data "github.com/coveo/terraform-auto-snippets/common_data"
	"github.com/coveo/terraform-auto-snippets/utils"
)

func getArgs(name string, uri url.URL, head *goquery.Selection) (arguments []data.Argument) {
	if head.Length() == 0 {
		return
	}

	var subArgument string
	_ = subArgument
	head.NextFilteredUntil("ul, p", "h2").Each(func(i int, section *goquery.Selection) {
		nodeType := section.Nodes[0].Data
		switch nodeType {
		case "ul":
			section.Find("li").Each(func(i int, arg *goquery.Selection) {
				const req = "(Required)"
				const opt = "(Optional)"
				argName, _ := arg.Find("a").First().Attr("name")
				description := strings.TrimSpace(strings.TrimPrefix(arg.Text(), argName+" - "))
				required := strings.Contains(description, req)
				if required {
					description = strings.Replace(strings.TrimSpace(description), req, "", 1)
				} else {
					description = strings.Replace(strings.TrimSpace(description), opt, "", 1)
				}

				url, _ := arg.Find("a[href]").Attr("href")
				arguments = append(arguments, data.Argument{
					Base: data.Base{
						Name:        argName,
						Description: utils.Trim(description),
						URL:         uri.String() + url,
					},
					Required: required,
				})
			})
		case "p":
		}
	})
	return
}
