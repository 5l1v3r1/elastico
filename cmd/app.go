package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"text/template"

	client "github.com/dutchcoders/elastico/client"
	json "github.com/dutchcoders/elastico/json"

	"github.com/codegangsta/cli"
	"github.com/op/go-logging"
)

func init() {
	cli.AppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}
USAGE:
   {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} {{if .Flags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}
   {{if .Version}}
VERSION:
   {{.Version}}
   {{end}}{{if len .Authors}}
AUTHOR(S):
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
COMMANDS:
   {{range .Commands}}{{join .Names ", "}}{{ "\t" }}{{.Usage}}
   {{end}}{{end}}{{if .Flags}}
GLOBAL OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}{{end}}{{if .Copyright }}
COPYRIGHT:
   {{.Copyright}}
   {{end}}
`

	cli.CommandHelpTemplate = `NAME:
   {{.HelpName}} - {{.Usage}}
USAGE:
   {{.HelpName}}{{if .Flags}}{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{if .Description}}
DESCRIPTION:
   {{.Description}}{{end}}{{if .Flags}}
OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}{{ end }}
`

	// The text template for the subcommand help topic.
	// cli.go uses text/template to render templates. You can
	// render custom help text by setting this variable.
	cli.SubcommandHelpTemplate = `NAME:
   {{.HelpName}} - {{.Usage}}
USAGE:
   {{.HelpName}} command{{if .Flags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
COMMANDS:
   {{range .Commands}}{{join .Names ", "}}{{ "\t" }}{{.Usage}}
   {{end}}{{if .Flags}}
OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}{{end}}
`
}

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

var e *client.Client

func run(fn func(*cli.Context) (json.M, error)) func(*cli.Context) error {
	return func(c *cli.Context) error {
		for _, flag := range c.Command.Flags {
			if _, ok := flag.(RequiredFlag); !ok {
				continue
			}

			if c.IsSet(flag.GetName()) {
				continue
			}

			err := fmt.Errorf("Missing parameter %s.", flag.GetName())

			ErrTemplate.Execute(os.Stderr, err)
			os.Exit(1)
		}

		resp, err := fn(c)
		if err != nil {
			ErrTemplate.Execute(os.Stderr, err)
			os.Exit(1)
		}
		if resp == nil {
			return nil
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

		return nil
	}
}

type App struct {
	*cli.App
}

func NewApp() *App {
	app := cli.NewApp()
	app.Name = "Elastico"
	app.Commands = []cli.Command{}

	app.Commands = append(app.Commands, searchCmds...)
	app.Commands = append(app.Commands, mappingCmds...)
	app.Commands = append(app.Commands, snapshotCmds...)
	app.Commands = append(app.Commands, indexCmds...)
	app.Commands = append(app.Commands, clusterCmds...)
	app.Commands = append(app.Commands, documentCmds...)
	app.Commands = append(app.Commands, templateCmds...)

	app.Version = "0.0.1"

	app.Before = func(context *cli.Context) error {
		backend1 := logging.NewLogBackend(os.Stderr, "", 0)
		backend1Leveled := logging.AddModuleLevel(backend1)
		backend1Leveled.SetLevel(logging.ERROR, "")
		logging.SetBackend(backend1Leveled)

		if context.GlobalBool("debug") {
			backend1Leveled.SetLevel(logging.DEBUG, "")
		}

		if t, err := client.New(context.String("host")); err != nil {
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
			Name:  "debug",
			Usage: "Enable debug mode",
		},
		cli.BoolFlag{
			Name:  "json",
			Usage: "Return json output",
		},
		cli.StringFlag{
			Name:   "host",
			Value:  "http://127.0.0.1:9200",
			EnvVar: "ELASTICO_HOST",
			Usage:  "Host to operate on",
		},
	}

	app.Flags = append(app.Flags, []cli.Flag{}...)

	return &App{
		app,
	}
}
