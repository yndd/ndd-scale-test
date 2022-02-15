package templ

import (
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/stoewer/go-strcase"
)

func ParseTemplate(templatefile string) (*template.Template, error) {

	fmt.Printf("templatefile: %s\n", templatefile)

	return template.Must(template.New(filepath.Base(templatefile)).Funcs(templateHelperFunctions).Funcs(sprig.TxtFuncMap()).ParseFiles(templatefile)), nil
	/*
		t := template.New(templatefile).Funcs(templateHelperFunctions).Funcs(sprig.TxtFuncMap())
		if _, err := t.ParseFiles(templatefile); err != nil {
			return nil, err
		}

		return t, nil
	*/
}

// templateHelperFunctions specifies a set of functions that are supplied as
// helpers to the templates that are used within this file.
var templateHelperFunctions = template.FuncMap{
	"inc":  func(i int) int { return i + 1 },
	"dec":  func(i int) int { return i - 1 },
	"mul":  func(p1 int, p2 int) int { return p1 * p2 },
	"mul3": func(p1, p2, p3 int) int { return p1 * p2 * p3 },
	"boolValue": func(b bool) int {
		if b {
			return 1
		} else {
			return 0
		}
	},
	"removeDashes": func(s string) string {
		return strings.ReplaceAll(s, "-", "")
	},
	"toUpperCamelCase": strcase.UpperCamelCase,
	"toLowerCamelCase": strcase.LowerCamelCase,
	"toKebabCase":      strcase.KebabCase,
	"toLower":          strings.ToLower,
	"toUpper":          strings.ToUpper,
}
