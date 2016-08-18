package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	json "github.com/dutchcoders/elastico/json"

	"github.com/codegangsta/cli"
)

var (
	_ = registerTemplate("document:put", `
{{- if .created }}Document created. 
{{- else -}}
Document updated (version {{ ._version }}). 
{{- end -}}`)
	_ = registerTemplate("document:delete", `
Document deleted. `)
	_ = registerTemplate("document:get", `{{._source | json}}`)
)

var documentCmds = []cli.Command{
	cli.Command{
		Name:        "document:get",
		Usage:       "Retrieve document from index",
		ArgsUsage:   "(-index {index}) (-type {type}) documentid",
		UsageText:   "elastico document:get ",
		Description: `Retrieve document from the index.`,
		Action:      run(runDocumentGet),
		BashComplete: func(context *cli.Context) {
			fmt.Println("BashCompl")
		},
		Flags: []cli.Flag{
			IndexFlag,
			TypeFlag,
		},
	},
	cli.Command{
		Name:        "document:delete",
		Usage:       "Delete document from index",
		Description: ``,
		Action:      run(runDocumentDelete),
		Flags: []cli.Flag{
			IndexFlag,
			TypeFlag,
		},
	},
	cli.Command{
		Name:        "document:put",
		Usage:       "Store document into index",
		Description: ``,
		Action:      run(runDocumentPut),
		Flags: []cli.Flag{
			IndexFlag,
			TypeFlag,
		},
	},
}

func runDocumentPut(c *cli.Context) (json.M, error) {
	if len(c.Args()) == 0 {
		return nil, fmt.Errorf("You need to supply the document id")
	}

	documentID := c.Args()[0]

	path := c.String("index")
	path = filepath.Join(path, c.String("type"))
	path = filepath.Join(path, documentID)

	var body interface{}
	if fi, err := os.Stdin.Stat(); err != nil {
		return nil, err
	} else if fi.Mode()&os.ModeNamedPipe > 0 {
		body = os.Stdin
	}

	req, err := e.NewRequest("PUT", documentID, body)
	if err != nil {
		return nil, err
	}

	var resp json.M
	if err := e.Do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func runDocumentGet(c *cli.Context) (json.M, error) {
	if len(c.Args()) == 0 {
		return nil, fmt.Errorf("You need to supply the document id")
	}

	documentID := c.Args()[0]

	path := c.String("index")
	path = filepath.Join(path, c.String("type"))
	path = filepath.Join(path, documentID)

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

func runDocumentDelete(c *cli.Context) (json.M, error) {
	if len(c.Args()) == 0 {
		return nil, fmt.Errorf("You need to supply the document id")
	}

	documentID := c.Args()[0]

	path := c.String("index")
	path = filepath.Join(path, c.String("type"))
	path = filepath.Join(path, documentID)

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
