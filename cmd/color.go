// Copyright © 2018 Shawn Catanzarite <me@shawncatz.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
