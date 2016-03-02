package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"text/template"

	elastico "elastico"
	json "elastico/json"

	"github.com/codegangsta/cli"
)

var (
	templates = map[string]*template.Template{}
)

func registerTemplate(cmd string, t string) *template.Template {
	var err error

	var usr *user.User
	if usr, err = user.Current(); err != nil {
		panic(err)
	}

	p := path.Join(usr.HomeDir, ".elastico")
	templatePath := filepath.Join(p, fmt.Sprintf("%s.template", cmd))

	if _, err := os.Stat(templatePath); err == nil {
		b, _ := ioutil.ReadFile(templatePath)
		templates[cmd] = template.Must(template.New("").Funcs(funcMap).Parse(string(b)))
	} else {

		templates[cmd] = template.Must(template.New("").Funcs(funcMap).Parse(t))
	}
	return templates[cmd]
}

var e *elastico.Elastico

func run(fn func(*cli.Context) (json.M, error)) func(*cli.Context) {
	return func(c *cli.Context) {
		resp, err := fn(c)
		if err != nil {
			// error template
			fmt.Println(err.Error())
			return
		}

		if resp == nil {
			return
		}

		if c.GlobalBool("json") {
			b, _ := json.MarshalIndent(resp, "", "  ")
			os.Stdout.Write(b)
		} else {
			if t, ok := templates[c.Command.Name]; ok {
				t.Execute(os.Stdout, resp)
			} else {
				b, _ := json.MarshalIndent(resp, "", "  ")
				os.Stdout.Write(b)
			}
		}
		return
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "Elastico"
	app.Commands = []cli.Command{}

	app.Commands = append(app.Commands, searchCmds...)
	app.Commands = append(app.Commands, snapshotCmds...)
	app.Commands = append(app.Commands, indexCmds...)
	app.Commands = append(app.Commands, clusterCmds...)

	app.Before = func(context *cli.Context) error {
		t, err := elastico.New(context.String("host"))
		if err != nil {
			panic(err)
		}
		e = t
		return nil
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name: "json",
		},
		cli.StringFlag{
			Name:   "host",
			Value:  "http://127.0.0.1:9200",
			EnvVar: "ELASTIC_HOST",
		},
	}

	app.Flags = append(app.Flags, []cli.Flag{}...)

	app.Run(os.Args)
}
