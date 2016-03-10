package main

import (
	"bytes"
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
	"github.com/op/go-logging"
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

var ErrTemplate = template.Must(template.New("").Funcs(funcMap).Parse(`Error: {{ .Error }}
`))

var e *elastico.Elastico

func run(fn func(*cli.Context) (json.M, error)) func(*cli.Context) {
	return func(c *cli.Context) {
		resp, err := fn(c)
		if err != nil {
			ErrTemplate.Execute(os.Stderr, err)
			os.Exit(1)
		}
		if resp == nil {
			return
		}

		buff := new(bytes.Buffer)
		defer func() {
			if err != nil {
				ErrTemplate.Execute(os.Stderr, err)
				os.Exit(1)
			}

			buff.WriteTo(os.Stdout)
		}()

		if c.GlobalBool("json") {
			b, _ := json.MarshalIndent(resp, "", "  ")
			buff.Write(b)
		} else {
			if t, ok := templates[c.Command.Name]; ok {
				err = t.Execute(buff, resp)
			} else {
				b, _ := json.MarshalIndent(resp, "", "  ")
				buff.Write(b)
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
		backend1 := logging.NewLogBackend(os.Stderr, "", 0)
		backend1Leveled := logging.AddModuleLevel(backend1)
		backend1Leveled.SetLevel(logging.ERROR, "")
		logging.SetBackend(backend1Leveled)

		if context.GlobalBool("debug") {
			backend1Leveled.SetLevel(logging.DEBUG, "")
		}

		if t, err := elastico.New(context.String("host")); err != nil {
			return err
		} else {
			e = t
			return nil
		}
	}

	app.EnableBashCompletion = true
	app.BashComplete = func(c *cli.Context) {
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name: "debug",
		},
		cli.BoolFlag{
			Name: "json",
		},
		cli.StringFlag{
			Name:   "host",
			Value:  "http://127.0.0.1:9200",
			EnvVar: "ELASTICO_HOST",
		},
		cli.StringFlag{
			Name:   "index",
			Value:  "",
			EnvVar: "ELASTICO_INDEX",
		},
		cli.StringFlag{
			Name:   "type",
			Value:  "",
			EnvVar: "ELASTICO_TYPE",
		},
	}

	app.Flags = append(app.Flags, []cli.Flag{}...)

	app.Run(os.Args)
}
