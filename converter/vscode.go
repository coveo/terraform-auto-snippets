package main

import (
	"encoding/json"
	"reflect"
	"strings"

	"io/ioutil"

	"github.com/kennygrant/sanitize"
)

type VscodeSnippet struct {
	Description string   `json:"description"`
	Prefix      string   `json:"prefix"`
	Body        []string `json:"body"`
}

func ProviderToVscodeSnippet(p Provider) *VscodeSnippet {
	prefix := "provider-" + sanitize.Name(p.Name)

	body := createSnippetBody(reflect.TypeOf(p).Name(), p.Name, p.URL, p.Arguments)
	return ToVscodeSnippet(p.Description, prefix, body)
}

func DataResourceToVscodeSnippet(d Dataresource, providerName string) *VscodeSnippet {
	prefix := "data-" + sanitize.Name(providerName) + "-" + sanitize.Name(d.Name)

	body := createSnippetBody(reflect.TypeOf(d).Name(), d.Name, d.URL, d.Arguments)
	return ToVscodeSnippet(d.Description, prefix, body)
}

func ResourceToVscodeSnippet(r Resource, providerName string) *VscodeSnippet {
	prefix := "res-" + sanitize.Name(providerName) + "-" + sanitize.Name(r.Name)

	body := createSnippetBody(reflect.TypeOf(r).Name(), r.Name, r.URL, r.Arguments)
	return ToVscodeSnippet(r.Description, prefix, body)
}

func ToVscodeSnippet(description string, prefix string, body []string) *VscodeSnippet {
	return &VscodeSnippet{
		Description: description,
		Prefix:      prefix,
		Body:        body,
	}
}

func createSnippetBody(kind string, name string, url string, args []Argument) []string {
	// // Get the kind of struct
	// kind := reflect.TypeOf(obj).Name()
	// Convert obj to standard Ressource

	//if we know the kind of resource convert to access flied
	// switch kind {
	// case "Provider":
	// 	p = obj.(Provider)
	// case "Dataresource ":
	// 	p := obj.(Dataresource)
	// case "Resource ":
	// 	p := obj.(Resource)
	// default:
	// 	p := obj.(TfObjectResource)
	// }
	// Get a clean name of resource (without space)
	cleanName := sanitize.Name(name)
	// Create the string for vscode
	body := []string{"# Configure the " + name,
		"# Doc : " + url,
		ResourceName(kind, cleanName) + " {"}

	// Add all args
	for _, a := range args {
		arg := "\t"
		// If the args is not required we comment it
		if a.Required == "false" {
			arg = arg + "#"
		}

		arg = arg + a.Name + " = "

		body = append(body, "\t# "+a.Description+" (require="+a.Required+")")
		body = append(body, arg)
	}

	body = append(body, "}")
	return body
}

func ResourceName(kind string, name string) string {
	var result string

	if kind == "Dataresource" {
		result = "data"
	} else {
		result = strings.ToLower(kind)
	}

	result = result + " \"" + name + "\""
	return result
}

// create snippet
func VscodeCreateSnippets(p map[string]Provider) {

	snippets := map[string]VscodeSnippet{}
	for _, v := range p {

		snippets[v.Name] = *ProviderToVscodeSnippet(v)
		for _, d := range v.Dataresources {
			snippets[v.Name+" "+d.Name] = *DataResourceToVscodeSnippet(d, v.Name)
		}

		for _, r := range v.Resources {
			snippets[v.Name+" "+r.Name] = *ResourceToVscodeSnippet(r, v.Name)
		}
		//result, _ := json.MarshalIndent(&snippet, "", " ")
	}

	result, _ := json.MarshalIndent(&snippets, "", " ")

	ioutil.WriteFile("../vscode/terraform-auto-snippets/snippets/snippets.json", result, 0644)
	//fmt.Print(string(result))
}
