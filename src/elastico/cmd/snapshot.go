package main

import (
	"fmt"
	"path/filepath"

	json "elastico/json"

	"github.com/codegangsta/cli"
)

var (
	_ = registerTemplate("snapshot:status", `=== Snapshots
{{range .snapshots }}
{{- .snapshot }} {{ .repository }} {{ .state }} {{.shards_stats.total}} {{.shards_stats.done}} {{.shards_stats.failed}} {{.stats.number_of_files}}{{.stats.time_in_millis}}{{.stats.total_size_in_bytes}} {{ .stats.processed_size_in_bytes }} {{ .stats.total_size_in_bytes }}
{{ end}}`)
)

var snapshotCmds = []cli.Command{
	cli.Command{
		Name:        "snapshot",
		Usage:       "",
		Description: ``,
		Action:      run(runSnapshotGet),
		Flags:       []cli.Flag{},
	},
	cli.Command{
		Name:        "snapshot:register",
		Usage:       "",
		Description: ``,
		Action:      run(runSnapshotRegister),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "type",
				Value: "fs",
			},
			cli.BoolFlag{
				Name: "compress",
			},
			Required(cli.StringFlag{
				Name: "location",
			}),
		},
	},
	cli.Command{
		Name:        "snapshot:execute",
		Usage:       "",
		Description: ``,
		Action:      run(runSnapshotExecute),
		Flags: []cli.Flag{
			Required(
				cli.StringFlag{
					Name: "repository",
				},
			),
			Required(
				cli.StringFlag{
					Name: "snapshot",
				},
			),
		},
	},
	cli.Command{
		Name:        "snapshot:status",
		Usage:       "",
		Description: ``,
		Action:      run(runSnapshotStatus),
		Flags: []cli.Flag{
			Required(
				cli.StringFlag{
					Name: "repository",
				},
			),
			Required(
				cli.StringFlag{
					Name: "snapshot",
				},
			),
		},
	},
}

func runSnapshotStatus(c *cli.Context) (json.M, error) {
	repository := c.String("repository")
	snapshot := c.String("snapshot")

	req, err := e.NewRequest("GET", fmt.Sprintf("/_snapshot/%s/%s/_status", repository, snapshot), nil)
	if err != nil {
		return nil, err
	}

	var resp json.M
	if err := e.Do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func runSnapshotExecute(c *cli.Context) (json.M, error) {
	repository := c.String("repository")
	snapshot := c.String("snapshot")

	body := json.M{}

	req, err := e.NewRequest("PUT", fmt.Sprintf("/_snapshot/%s/%s?wait_for_completion=true", repository, snapshot), body)
	if err != nil {
		return nil, err
	}

	var resp json.M
	if err := e.Do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func runSnapshotRegister(c *cli.Context) (json.M, error) {
	if len(c.Args()) == 0 {
		return nil, fmt.Errorf("You need to supply the name of the snapshot")
	}

	name := c.Args()[0]
	location := c.String("location")
	type_ := c.String("type")

	body := json.M{
		"settings": json.M{
			"compress": "true",
			"location": location,
		},
		"type": type_,
	}

	req, err := e.NewRequest("PUT", fmt.Sprintf("/_snapshot/%s", name), body)
	if err != nil {
		return nil, err
	}

	var resp json.M
	if err := e.Do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func runSnapshotGet(c *cli.Context) (json.M, error) {
	req, err := e.NewRequest("GET", filepath.Join("/_snapshot/_all"), nil)
	if err != nil {
		return nil, err
	}

	var resp json.M
	if err := e.Do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}
