package main

import (
	"fmt"
	"os"
	"path/filepath"

	json "elastico/json"

	"github.com/codegangsta/cli"
)

var (
	_ = registerTemplate("template:create", `Template created.
`)
	_ = registerTemplate("template:list", `== Templates
Name				      Indexes	                         Order
===================================== ========================= ==============
{{range $name, $template := . }}
{{- $name | yellow | printf "%-46s" }} {{ $template.template | printf "%-20s" }}{{ $template.order| printf "%20.0f" }}
{{ end}}`)
	_ = registerTemplate("template:delete", `Template deleted.
`)
)

var templateCmds = []cli.Command{
	cli.Command{
		Name:        "template:delete",
		Usage:       "Delete template",
		Description: ``,
		Action:      run(runTemplateDelete),
		Flags: []cli.Flag{
			TemplateRequiredFlag,
		},
	},
	cli.Command{
		Name:        "template:list",
		Usage:       "List templates",
		Description: ``,
		Action:      run(runTemplateList),
		Flags:       []cli.Flag{},
	},
	cli.Command{
		Name:        "template:create",
		Usage:       "Create template",
		Description: ``,
		Action:      run(runTemplateCreate),
		Flags: []cli.Flag{
			TemplateRequiredFlag,
		},
	},
}

func runTemplateDelete(c *cli.Context) (json.M, error) {
	template := c.String("template")

	path := filepath.Join("_template", template)

	req, err := e.NewRequest("DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	var resp json.M
	if err := e.Do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func runTemplateList(c *cli.Context) (json.M, error) {
	path := filepath.Join("_template")

	req, err := e.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var resp json.M
	if err := e.Do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}
func runTemplateCreate(c *cli.Context) (json.M, error) {
	template := c.String("template")

	path := filepath.Join("_template", template)

	var body interface{}
	if fi, err := os.Stdin.Stat(); err != nil {
		return nil, fmt.Errorf("No stdin body template")
	} else if fi.Mode()&os.ModeNamedPipe > 0 {
		body = os.Stdin
	}

	req, err := e.NewRequest("PUT", path, body)
	if err != nil {
		return nil, err
	}

	var resp json.M
	if err := e.Do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}
