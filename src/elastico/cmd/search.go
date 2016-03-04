package main

import (
	"fmt"
	"os"
	"path/filepath"

	json "elastico/json"

	"github.com/codegangsta/cli"
)

var (
	_ = registerTemplate("search", `== Search
Took:	    {{ .took | time}}
Total:	    {{ .hits.total}}
Max score:  {{ .hits.max_score}}

{{ if gt (.hits.hits | len) 0 -}}
== Hits
Index				   Type	                  ID		                           Score
================================== ====================== ======================================== ====================
{{range $hit := .hits.hits -}}
{{- $hit._index | yellow | printf "%-44s" }}{{ $hit._type | printf "%-22s" }} {{ $hit._id | printf "%-40s" }} {{ $hit._score |  printf "%20.2f" }}
{{end -}}
{{- else -}}
{{end -}}
{{ if gt (.aggregations | len) 0 }}
== Aggregations
{{range $name, $agg := .aggregations -}}
{{- $hasBuckets := (has $agg "buckets") -}}
{{ if $hasBuckets -}}
+ {{ $name }}
{{range $bucket := $agg.buckets -}}
{{ "  " }}{{ $bucket.key }}: {{ $bucket.doc_count }}
{{- end }}
  other: {{ $agg.sum_other_doc_count }}
{{- else -}}
+ {{ $name }}: {{ $agg.value -}} 
{{end }}
{{ end }}
{{end -}}
`)
	_ = registerTemplate("put", `{{. | json}}`)
	_ = registerTemplate("get", `{{._source | json}}`)
)

var searchCmds = []cli.Command{
	cli.Command{
		Name:        "get",
		Usage:       "",
		Description: ``,
		Action:      run(runGet),
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
	cli.Command{
		Name:        "put",
		Usage:       "",
		Description: ``,
		Action:      run(runPut),
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

func runPut(c *cli.Context) (json.M, error) {
	index := c.String("index")
	type_ := c.String("type")

	path := c.Args()[0]
	if type_ != "" {
		path = filepath.Join(type_, path)
	}
	if index != "" {
		path = filepath.Join(index, path)
	}

	var body interface{}
	if fi, err := os.Stdin.Stat(); err != nil {
		return nil, err
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

func runGet(c *cli.Context) (json.M, error) {
	index := c.String("index")
	type_ := c.String("type")

	path := c.Args()[0]
	if type_ != "" {
		path = filepath.Join(type_, path)
	}
	if index != "" {
		path = filepath.Join(index, path)
	}

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
