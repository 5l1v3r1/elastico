package console

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var DefaultConsole = New()

func Scanln(ask string, args ...interface{}) (string, error) {
	return DefaultConsole.Scanln(ask, args...)
}

func New() *Console {
	return &Console{}
}

type Console struct {
}

func (c *Console) Scanln(ask string, args ...interface{}) (string, error) {
	r := bufio.NewReader(os.Stdin)
	fmt.Printf(ask, args...)

	input, err := r.ReadString('\n')
	return strings.TrimSpace(input), err
}
