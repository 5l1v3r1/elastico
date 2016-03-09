package main

import (
	"elastico/json"
	"fmt"
	"regexp"
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
		"trim": func(v string) string {
			re := regexp.MustCompile("\\s+")
			return re.ReplaceAllLiteralString(v, " ")
		},
		"pretty": pretty.Print,
		"json": func(v interface{}) (string, error) {
			b, err := json.MarshalIndent(v, "", "  ")
			return string(b), err
		},
		"time": func(v interface{}) (string, error) {
			return fmt.Sprintf("%d", (uint64(v.(float64)))/1000), nil
			// return humanize.Time(uint64(v.(float64))), nil
		},
		"bytes": func(v interface{}) (string, error) {
			return humanize.Bytes(uint64(v.(float64))), nil
		},
		"has": func(v interface{}, key string) (bool, error) {
			if m, ok := v.(json.M); !ok {
				return false, nil
			} else {
				_, match := m[key]
				return match, nil
			}
		},
	}
)

func Pad(v string, count int) string {
	return v // + strconv.Itoa(arg)
}
