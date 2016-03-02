package main

import (
	"fmt"
	"text/template"

	humanize "github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/kr/pretty"
)

var (
	funcMap = template.FuncMap{
		"pad":    Pad,
		"red":    color.RedString,
		"green":  color.GreenString,
		"yellow": color.YellowString,
		"blue":   color.BlueString,
		"pretty": pretty.Print,
		"time": func(v interface{}) (string, error) {
			return fmt.Sprintf("%d", (uint64(v.(float64)))/1000), nil
			// return humanize.Time(uint64(v.(float64))), nil
		},
		"bytes": func(v interface{}) (string, error) {
			return humanize.Bytes(uint64(v.(float64))), nil
		},
	}
)

func Pad(v string, count int) string {
	return v // + strconv.Itoa(arg)
}
