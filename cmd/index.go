package cmd

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	client "github.com/dutchcoders/elastico/client"
	json "github.com/dutchcoders/elastico/json"

	pb "github.com/cheggaaa/pb"

	"github.com/codegangsta/cli"
)

var (
	_ = registerTemplate("index:open", `Index opened.`)
	_ = registerTemplate("index:close", `Index closed.`)
	_ = registerTemplate("index:create", `Index created.`)
	_ = registerTemplate("index:delete", `Index deleted.`)
	_ = registerTemplate("index:stats", `== Shards
failed		    {{ ._shards.failed }}
successful	    {{ ._shards.successful }}
total		    {{ ._shards.total }}

== Indices
Name				  	            Count	       Deleted		  Querytime	  Size in bytes
================================== ====================== ==================== ==================== ===================
{{range $name, $index := .indices }}
{{- $name | yellow | printf "%-46s" }}{{ $index.primaries.docs.count | printf "%20.0f" }} {{ $index.primaries.docs.deleted | printf "%20.0f" }} {{ $index.primaries.search.query_time_in_millis | time | printf "%20s" }}{{ $index.primaries.store.size_in_bytes | bytes | printf "%20s" }}
{{ end}}`)

	_ = registerTemplate("index:recovery", `== Shards
failed		    {{ .shards.failed }}
successful	    {{ .shards.successful }}
total		    {{ .shards.total }}

== Indices
Name				  	          NumDocs		MaxDoc		 DeletedDoc	  Size in bytes
================================== ====================== ==================== ==================== ===================
{{range $name, $index := . }}
{{- $name | yellow | printf "%-46s" }}{{ $index.docs.num_docs | printf "%20.0f" }} {{ $index.docs.max_doc | printf "%20.0f" }} {{ $index.docs.deleted_docs | printf "%20.0f" }}{{ $index.index.size_in_bytes | printf "%20.0f" }}
{{ end}}`)

	_ = registerTemplate("index:get", `=== Index
{{range $name, $index := . }}
{{- $name | printf "== %s"}}
{{ if .aliases | len | gt 0 -}}
Aliases:
{{range $alias := .aliases -}}
{{ $alias }}
{{- end }}
{{- end -}}
Creation date:	    {{ $index.settings.index.creation_date }}
Replicas:	    {{ $index.settings.index.number_of_replicas }}
Shards:		    {{ $index.settings.index.number_of_shards }}
UUID:		    {{ $index.settings.index.uuid }}
Version created:    {{ $index.settings.index.version.created }}

{{ end}}
`)
)

var indexCmds = []cli.Command{
	cli.Command{
		Name:        "index:create",
		Usage:       "Create index",
		Description: ``,
		Action:      run(runIndexCreate),
		Flags: []cli.Flag{
			IndexRequiredFlag,
			cli.IntFlag{
				Name:  "replicas",
				Value: 1,
			},
			cli.IntFlag{
				Name:  "shards",
				Value: 5,
			},
		},
	},
	cli.Command{
		Name:        "index:update",
		Usage:       "Update index settings",
		Description: ``,
		Action:      run(runIndexUpdate),
		Flags: []cli.Flag{
			IndexRequiredFlag,
			cli.IntFlag{
				Name: "replicas",
			},
			cli.StringFlag{
				Name: "refresh-interval",
			},
		},
	},
	cli.Command{
		Name:        "index:get",
		Aliases:     []string{"index"},
		Usage:       "Retrieve index info",
		Description: ``,
		Action:      run(runIndexGet),
		Flags: []cli.Flag{
			IndexFlag,
		},
	},
	cli.Command{
		Name:        "index:delete",
		Usage:       "Delete index",
		Description: ``,
		Action:      run(runIndexDelete),
		Flags: []cli.Flag{
			IndexRequiredFlag,
		},
	},
	cli.Command{
		Name:        "index:copy",
		Usage:       "Copy index to other host or index",
		Description: ``,
		Action:      run(runIndexCopy),
		Flags: []cli.Flag{
			IndexRequiredFlag,
			TypeFlag,
			cli.StringFlag{
				Name: "from",
			},
		},
	},
	cli.Command{
		Name:        "index:recovery",
		Usage:       "Retrieve index recovery insights",
		Description: ``,
		Action:      run(runIndexRecovery),
		Flags: []cli.Flag{
			IndexRequiredFlag,
		},
	},
	cli.Command{
		Name:        "index:stats",
		Usage:       "Get index statistics",
		Description: ``,
		Action:      run(runIndexStats),
		Flags: []cli.Flag{
			IndexFlag,
		},
	},
	cli.Command{
		Name:        "index:open",
		Usage:       "Open index",
		Description: ``,
		Action:      run(runIndexOpen),
		Flags: []cli.Flag{
			IndexRequiredFlag,
		},
	},
	cli.Command{
		Name:        "index:close",
		Usage:       "Close index",
		Description: ``,
		Action:      run(runIndexClose),
		Flags: []cli.Flag{
			IndexRequiredFlag,
		},
	},
}

func runIndexOpen(c *cli.Context) (json.M, error) {
	path := c.String("index")
	path = filepath.Join(path, "_open")

	req, err := e.NewRequest("POST", path, nil)
	if err != nil {
		return nil, err
	}

	var resp json.M
	if err := e.Do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func runIndexClose(c *cli.Context) (json.M, error) {
	path := c.String("index")
	path = filepath.Join(path, "_close")

	req, err := e.NewRequest("POST", path, nil)
	if err != nil {
		return nil, err
	}

	var resp json.M
	if err := e.Do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func runIndexUpdate(c *cli.Context) (json.M, error) {
	path := c.String("index")
	path = filepath.Join(path, "_settings")

	body := json.M{
		"index": json.M{},
	}

	if c.IsSet("replicas") {
		body["index"].(json.M).Set("number_of_replicas", c.Int("replicas"))
	}

	if c.IsSet("refresh-interval") {
		body["index"].(json.M)["refresh_interval"] = c.String("refresh-interval")
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

func runIndexGet(c *cli.Context) (json.M, error) {
	path := c.String("index")

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

func runIndexCreate(c *cli.Context) (json.M, error) {
	path := c.String("index")

	var body interface{}
	if fi, err := os.Stdin.Stat(); err != nil {
		return nil, err
	} else if fi.Mode()&os.ModeNamedPipe > 0 {
		body = os.Stdin
	} else {
		body = json.M{
			"settings": json.M{
				"replicas": c.Int("replicas"),
				"shards":   c.Int("shards"),
			},
		}
	}

	req, err := e.NewRequest("POST", path, body)
	if err != nil {
		return nil, err
	}

	var resp json.M
	if err := e.Do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func runIndexDelete(c *cli.Context) (json.M, error) {
	path := c.String("index")

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

func runIndexStats(c *cli.Context) (json.M, error) {
	path := c.String("index")
	path = filepath.Join(path, "_stats")

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

func runIndexRecovery(c *cli.Context) (json.M, error) {
	path := c.String("index")
	path = filepath.Join(path, "_recovery")

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

func runIndexCopy(c *cli.Context) (json.M, error) {
	path := c.String("index")
	path = filepath.Join(path, fmt.Sprintf("_search?sort=_doc&scroll=1m&size=1"))

	if len(c.Args()) == 0 {
		return nil, fmt.Errorf("You need to supply the destination location")
	}

	rel, err := url.Parse(c.Args()[0])
	if err != nil {
		return nil, err
	}

	dst, err := client.New(e.BaseURL.ResolveReference(rel).String())
	if err != nil {
		return nil, err
	}

	var body interface{}
	if fi, err := os.Stdin.Stat(); err != nil {
		return nil, err
	} else if fi.Mode()&os.ModeNamedPipe > 0 {
		body = os.Stdin
	}

	var bar *pb.ProgressBar
	defer func() {
		if bar == nil {
			return
		}

		if err == nil {
			bar.FinishPrint("Done!")
		}
	}()

	count := float64(0)

	req, err := e.NewRequest("GET", path, body)
	for err == nil {
		var resp json.M
		if err = e.Do(req, &resp); err != nil {
			break
		}

		hits := resp["hits"].(json.M)["hits"].([]interface{})
		if len(hits) == 0 {
			break
		}

		total := resp["hits"].(json.M)["total"].(float64)
		if bar == nil {
			bar = pb.StartNew(int(total))
		}

		var buffer bytes.Buffer

		enc := json.NewEncoder(&buffer)
		for _, val := range hits {
			hit := val.(json.M)
			if err = enc.Encode(json.M{
				"index": json.M{
					"_id": hit["_id"],
				},
			}); err != nil {
				break
			}

			// should not decode and encode, just bytes
			if err = enc.Encode(hit["_source"]); err != nil {
				break
			}
		}
		if err != nil {
			break
		}

		req, err = dst.NewRequest("POST", "_bulk", &buffer)
		if err != nil {
			break
		}

		err = dst.Do(req, &resp)
		if err != nil {
			break
		}

		count += float64(len(hits))

		bar.Set(int(count))

		scrollID := resp["_scroll_id"]
		req, err = e.NewRequest("GET", "/_search/scroll", json.M{
			"scroll":    "1m",
			"scroll_id": scrollID,
		})
	}

	return nil, err
}
