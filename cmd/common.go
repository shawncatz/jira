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
	"fmt"
	"io/ioutil"

	"github.com/andygrunwald/go-jira"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/viper"
)

var a aurora.Aurora

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

func printErrResponse(response *jira.Response) {
	r := response.Response
	printErr("Jira error: %d : %s\n", r.StatusCode, r.Status)
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	if len(b) > 0 {
		printErr("%s\n", b)
	}
}

func issueURL(id string) string {
	return viper.GetString("jira.base") + "/browse/" + id
}

func getBoards() (list []jira.Board, err error) {
	project := viper.GetString("jira.project")
	options := &jira.BoardListOptions{ProjectKeyOrID: project}
	br, response, err := jiraClient.Board.GetAllBoards(options)
	if err != nil {
		printErr("Error: %s", err)
		if debug {
			printErrResponse(response)
		}
		return list, err
	}
	return br.Values, err
}

func getSprints(id int) (list []jira.Sprint, err error) {
	sr, response, err := jiraClient.Board.GetAllSprints(fmt.Sprintf("%d", id))
	if err != nil {
		printErr("Error: %s", err)
		if debug {
			printErrResponse(response)
		}
		return list, err
	}
	return sr, err
}

// if we need to walk the response, something like this
//func getSprintsPage(boardId, start, limit int) ([]jira.Sprint, int, int, error) {
//	sr, response, err := jiraClient.Board.GetAllSprints(fmt.Sprintf("%d", boardId))
//
//}
