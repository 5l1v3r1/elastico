package main

import (
	"os"
	"path/filepath"

	json "elastico/json"

	"github.com/codegangsta/cli"
)

var (
	_ = registerTemplate("search", `== Search
	{{ . | pretty }}
    `)
)

var searchCmds = []cli.Command{
	cli.Command{
		Name:        "search",
		Usage:       "",
		Description: ``,
		Action:      run(runSearch),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "index",
				Value: "",
			},
			cli.StringFlag{
				Name:  "type",
				Value: "",
			},
		},
	},
}

func runSearch(c *cli.Context) (json.M, error) {
	index := c.String("index")
	type_ := c.String("type")

	path := "_search"
	if type_ != "" {
		path = filepath.Join(type_, path)
	}
	if index != "" {
		path = filepath.Join(index, path)
	}

	var body interface{}
	if len(c.Args()) > 0 {
		body = json.M{
			"query": json.M{
				"simple_query_string": json.M{
					"query": c.Args()[0],
				},
			},
		}
	} else if fi, err := os.Stdin.Stat(); err != nil {
	} else if fi.Mode()&os.ModeNamedPipe > 0 {
		body = os.Stdin
	}

	req, err := e.NewRequest("GET", path, body)
	if err != nil {
		return nil, err
	}

	var resp json.M
	if err := e.Do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}
