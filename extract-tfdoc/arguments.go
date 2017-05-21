package main

import (
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"strings"

	"fmt"
	data "github.com/coveo/terraform-auto-snippets/common_data"
	"github.com/coveo/terraform-auto-snippets/utils"
)

func getArgs(name string, uri url.URL, head *goquery.Selection) (arguments []data.Argument) {
	if head.Length() == 0 {
		return
	}

	var tempArguments []*data.Argument
	argsMap := make(map[string][]*data.Argument)
	var subArgs []string
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
				newArg := data.Argument{
					Base: data.Base{
						Name:        argName,
						Description: utils.Trim(description),
						URL:         uri.String() + url,
					},
					Required: required,
				}
				argsMap[argName] = append(argsMap[argName], &newArg)
				if len(subArgs) == 0 {
					tempArguments = append(tempArguments, &newArg)
				} else {
					for _, subArg := range subArgs {
						for _, arg := range argsMap[subArg] {
							arg.Fields = append(arg.Fields, newArg)
						}
					}
				}
			})
		case "p":
			if i == 0 {
				// We ignore the p element if it comes first
				return
			}

			subArgs = make([]string, 0)
			section.Find("code").Each(func(i int, block *goquery.Selection) {
				argRef := block.Text()
				if block.Length() == 1 {
					if _, ok := argsMap[argRef]; !ok {
						argRef = fmt.Sprintf("Black hole %s", argRef)
						// utils.PrintWarning("%s for %s", argRef, name)
						argsMap[argRef] = []*data.Argument{&data.Argument{}}
					}
					subArgs = append(subArgs, argRef)
				}
			})
		}
	})

	arguments = make([]data.Argument, len(tempArguments))
	for i := range tempArguments {
		arguments[i] = *tempArguments[i]
	}
	return
}
