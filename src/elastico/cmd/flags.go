package main

import "github.com/codegangsta/cli"

var IndexFlag = cli.StringFlag{
	Name:  "index",
	Value: "_all",
}

type RequiredFlag struct {
	cli.Flag
}

func (rsf RequiredFlag) Required() bool {
	return true
}

func Required(flag cli.Flag) cli.Flag {
	return RequiredFlag{
		flag,
	}
}

var (
	IndexRequiredFlag = RequiredFlag{
		cli.StringFlag{
			Name:   "index",
			Value:  "",
			EnvVar: "ELASTICO_INDEX",
		},
	}
	TypeRequiredFlag = RequiredFlag{
		cli.StringFlag{
			Name:   "type",
			Value:  "",
			EnvVar: "ELASTICO_TYPE",
		},
	}
	TemplateRequiredFlag = RequiredFlag{
		cli.StringFlag{
			Name:   "template",
			Value:  "",
			EnvVar: "",
		},
	}
	FieldRequiredFlag = RequiredFlag{
		cli.StringFlag{
			Name:  "field",
			Value: "",
		},
	}
	TypeFlag = cli.StringFlag{
		Name:  "type",
		Value: "",
	}

	FieldFlag = cli.StringFlag{
		Name:  "field",
		Value: "",
	}
)
