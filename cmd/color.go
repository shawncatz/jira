package cmd

import (
	"regexp"

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
