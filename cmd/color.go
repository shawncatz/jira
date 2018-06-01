package cmd

import (
	"regexp"

	"fmt"
	"github.com/logrusorgru/aurora"
)

var a aurora.Aurora
var re = regexp.MustCompile("\n")

func init() {
	a = aurora.NewAurora(true)
}

func white(v string) aurora.Value {
	return a.Bold(v).Gray()
}

func gray(v string) aurora.Value {
	return a.Gray(v)
}

func cyan(v string) aurora.Value {
	return a.Bold(v).Cyan()
}

func red(v string) aurora.Value {
	return a.Bold(v).Red()
}

func printErr(format string, a ...interface{}) (n int, err error) {
	return fmt.Println(red(fmt.Sprintf(format, a...)))
}
