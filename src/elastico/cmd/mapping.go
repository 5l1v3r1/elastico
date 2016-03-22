package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	console "elastico/console"
	json "elastico/json"

	"github.com/codegangsta/cli"
)

var (
	_ = registerTemplate("mapping:get", `{{ . | json }}
`)
	_ = registerTemplate("mapping:put", `Mapping updated.
`)
	_ = registerTemplate("mapping:edit", `Mapping updated.
`)
	// should we show original values here? str [start: end]
	_ = registerTemplate("analyze", `=== Analyze
{{ range .tokens }}
Offset:		{{ .start_offset }} - {{ .end_offset }}
Position:	{{ .position }}
Token:		{{ .token | yellow }}
Type:		{{ .type  }}
{{ end }}
`)
)

var mappingCmds = []cli.Command{
	cli.Command{
		Name:        "mapping:get",
		Usage:       "Get index mappings",
		Description: ``,
		Action:      run(runMappingGet),
		Flags: []cli.Flag{
			IndexRequiredFlag,
			TypeFlag,
		},
	},
	cli.Command{
		Name:        "mapping:put",
		Usage:       "Put index mappings",
		Description: ``,
		Action:      run(runMappingPut),
		Flags: []cli.Flag{
			IndexRequiredFlag,
			TypeFlag,
		},
	},
	cli.Command{
		Name:        "mapping:edit",
		Usage:       "Edit index mappings using editor",
		Description: ``,
		Action:      run(runMappingEdit),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "editor",
				EnvVar: "EDITOR",
				Value:  "vim",
			},
			IndexRequiredFlag,
			TypeFlag,
		},
	},
	cli.Command{
		Name:        "analyze",
		Usage:       "Analyze arguments by index and field",
		Description: ``,
		Action:      run(runAnalyze),
		Flags: []cli.Flag{
			IndexRequiredFlag,
			FieldRequiredFlag,
		},
	},
}

func runAnalyze(c *cli.Context) (json.M, error) {
	if len(c.Args()) == 0 {
		return nil, fmt.Errorf("You need to supply the text to analyze as argument")
	}

	path := c.String("index")
	path = filepath.Join(path, fmt.Sprintf("_analyze?field=%s", c.String("field")))

	body := bytes.NewBufferString(strings.Join(c.Args(), " "))
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

func runMappingGet(c *cli.Context) (json.M, error) {
	index := c.String("index")
	type_ := c.String("type")

	path := index
	path = filepath.Join(path, "_mapping", type_)

	req, err := e.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var resp json.M
	if err := e.Do(req, &resp); err != nil {
		return nil, err
	}

	if type_ != "" {
		return resp[index].(json.M)["mappings"].(json.M)[type_].(json.M), nil
	}

	return resp[index].(json.M)["mappings"].(json.M), nil
}

func runMappingPut(c *cli.Context) (json.M, error) {
	index := c.String("index")
	type_ := c.String("type")

	path := index
	if type_ != "" {
		path = filepath.Join(path, "_mapping", type_)
	}

	var body interface{}
	if fi, err := os.Stdin.Stat(); err != nil {
		return nil, fmt.Errorf("No stdin body mapping")
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

func runMappingEdit(c *cli.Context) (json.M, error) {
	index := c.String("index")
	type_ := c.String("type")

	path := index
	if type_ != "" {
		path = filepath.Join(path, "_mapping", type_)
	}

	var body interface{}
	req, err := e.NewRequest("GET", path, body)
	if err != nil {
		return nil, err
	}

	var resp json.M
	if err := e.Do(req, &resp); err != nil {
		return nil, err
	}

	if type_ != "" {
		resp = resp[index].(json.M)["mappings"].(json.M)
	} else {
		resp = json.M{
			"mappings": resp[index].(json.M)["mappings"].(json.M),
		}
	}

	tf, err := ioutil.TempFile("", "elastico-")
	if err != nil {
		return nil, err
	}

	defer os.Remove(tf.Name())

	b, _ := json.MarshalIndent(resp, "", "  ")
	tf.Write(b)
	tf.Close()

	for {
		cmd := exec.Command(c.String("editor"), tf.Name())
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Start(); err != nil {
			return nil, err
		}

		if err := cmd.Wait(); err != nil {
			return nil, err
		}

		if tf, err = os.Open(tf.Name()); err != nil {
			return nil, err
		}

		if err = json.NewDecoder(tf).Decode(&body); err != nil {
			input, _ := console.Scanln("Could not parse json file: %s. Continue [a,r]?", err.Error())
			if input == "a" {
				return nil, err
			}
			continue
		}

		req, err = e.NewRequest("PUT", path, body)
		if err != nil {
			return nil, err
		}

		if err := e.Do(req, &resp); err != nil {
			input, _ := console.Scanln("Error while updating mapping: %s. Continue [a,r] ?", err.Error())
			if input == "a" {
				return nil, err
			}
			continue
		} else {
			break
		}
	}

	return resp, nil
}
