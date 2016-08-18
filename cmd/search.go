package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	json "github.com/dutchcoders/elastico/json"

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
{{range $name, $highlights := $hit.highlight -}}
{{ $name | green }}: {{ range $highlight := $highlights -}} {{ $highlight | trim }} {{end }}
{{end -}}
{{end -}}
{{- else -}}
{{end -}}
{{ if (has . "aggregations") }}
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
{{end -}}
`)
)

var searchCmds = []cli.Command{
	cli.Command{
		Name:        "search",
		Usage:       "Search through index (and type)",
		Description: ``,
		Action:      run(runSearch),
		Flags: []cli.Flag{
			IndexFlag,
			TypeFlag,
			FieldFlag,
			cli.IntFlag{
				Name:  "size",
				Value: 10,
			},
			cli.BoolFlag{
				Name: "disable-highlight",
			},
		},
	},
}

// highlight
func runSearch(c *cli.Context) (json.M, error) {
	path := c.String("index")
	path = filepath.Join(path, c.String("type"))
	path = filepath.Join(path, fmt.Sprintf("_search?size=%d", c.Int("size")))

	var body interface{}
	if len(c.Args()) > 0 {
		body = json.M{
			"query": json.M{
				"simple_query_string": json.M{
					"query": c.Args()[0],
				},
			},
		}

		if !c.Bool("disable-highlight") {
			body.(json.M).Set("highlight", json.M{
				"require_field_match": false,
				"pre_tags":            []string{"\x1b[93m"},
				"post_tags":           []string{"\x1b[0m"},
				"tag_schema":          "styled",
				"fields": json.M{
					"*": json.M{},
				},
			})
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
