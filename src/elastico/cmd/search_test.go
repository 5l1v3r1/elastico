package main

import (
	"flag"
	"testing"

	elastico "elastico"

	"github.com/codegangsta/cli"
)

func TestSearch(t *testing.T) {
	//	setup("test-index", "test-type")
	app := cli.NewApp()
	if val, err := elastico.New("http://127.0.0.1:9200/"); err != nil {
		t.Fatal(err)
	} else {
		e = val
	}

	globalSet := flag.NewFlagSet("test", 0)
	globalSet.Set("index", "remco2")
	globalSet.Set("type", "remco2")
	globalCtx := cli.NewContext(app, globalSet, nil)

	set := flag.NewFlagSet("test", 0)
	set.Parse([]string{"-size", "10"})

	c := cli.NewContext(nil, set, globalCtx)

	run(runSearch)(c)
}
