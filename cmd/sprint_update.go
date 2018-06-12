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

	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
)

// sprintUpdateCmd represents the sprintUpdate command
var sprintUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update list of sprints in configuration",
	Long:  "Update list of sprints in configuration",
	Run: func(cmd *cobra.Command, args []string) {
		answer, err := askBoard()
		if err != nil {
			printErr("Error: %s\n", err)
			return
		}
		sprints, err := getSprints(answer, false)

		list := []string{}
		for _, s := range sprints {
			list = append(list, s.Name)
		}

		cfg.Jira.Sprints = list
		cfg.Save()
	},
}

func init() {
	sprintCmd.AddCommand(sprintUpdateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sprintUpdateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sprintUpdateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func askBoard() (int, error) {
	boards, err := getBoards()
	if err != nil {
		printErr("error: %s", err)
		return 0, err
	}

	if len(boards) == 0 {
		return 0, fmt.Errorf("no boards found")
	}

	if len(boards) == 1 {
		return boards[0].ID, nil
	}

	boardsList := []string{}
	boardsMap := map[string]int{}
	for _, b := range boards {
		boardsList = append(boardsList, b.Name)
		boardsMap[b.Name] = b.ID
	}
	fmt.Printf("%#v\n%#v\n", boardsList, boardsMap)

	answer := ""
	q := &survey.Select{
		Message: "Choose a board:",
		Options: boardsList,
	}
	survey.AskOne(q, &answer, survey.Required)

	return boardsMap[answer], nil
}
