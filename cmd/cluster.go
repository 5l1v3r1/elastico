package cmd

import (
	json "github.com/dutchcoders/elastico/json"

	"github.com/codegangsta/cli"
)

var (
	_ = registerTemplate("cluster:state", `=== Cluster state
cluster name		    {{ .cluster_name }}
master node		    {{ .master_node }}

= Nodes
{{range $name, $node := .nodes -}}
{{ $name }}		    {{ $node.name }}		    {{ $node.transport_address }}  
{{end}}
`)
	_ = registerTemplate("cluster:health", `=== Cluster health
name			    {{ .cluster_name }}
{{- if eq .status "green" }}
status			    {{ .status | green }}
{{ end -}}
{{- if eq .status "yellow" }}
status			    {{ .status | yellow }}
{{ end -}}
{{- if eq .status "red" }}
status			    {{ .status | red }}
{{ end -}}
# nodes			    {{ .number_of_nodes }}
# data nodes		    {{ .number_of_data_nodes }}
# active shards		    {{ .active_shards }}
# active primary shards	    {{ .active_primary_shards }}
# initializing shards	    {{ .initializing_shards }}
# relocating shards	    {{ .relocating_shards }}
# unassigned shards	    {{ .unassigned_shards }}
`)
)

var clusterCmds = []cli.Command{
	cli.Command{
		Name:        "cluster:state",
		Usage:       "Retrieve cluster state",
		Description: ``,
		Action:      run(runClusterState),
		Flags:       []cli.Flag{},
	},
	cli.Command{
		Name:        "cluster:health",
		Usage:       "Retrieve cluster health",
		Description: ``,
		Action:      run(runClusterHealth),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "type",
				Value: "fs",
			},
			cli.BoolFlag{
				Name: "compress",
			},
			cli.StringFlag{
				Name: "location",
			},
		},
	},
}

func runClusterState(c *cli.Context) (json.M, error) {
	req, err := e.NewRequest("GET", "/_cluster/state", nil)
	if err != nil {
		return nil, err
	}

	var resp json.M
	if err := e.Do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func runClusterHealth(c *cli.Context) (json.M, error) {
	req, err := e.NewRequest("GET", "/_cluster/health", nil)
	if err != nil {
		return nil, err
	}

	var resp json.M
	if err := e.Do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}
