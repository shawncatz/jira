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
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update configuration from Jira",
	Run: func(cmd *cobra.Command, args []string) {
		board, err := askBoard()
		if err != nil {
			printErr("Error: %s\n", err)
			return
		}

		sprints, err := getSprints(board.ID, false)

		list := []*JiraSprint{}
		for _, s := range sprints {
			list = append(list, &JiraSprint{ID: s.ID, Name: s.Name})
		}

		cfg.Jira.Board = &JiraBoard{ID: board.ID, Name: board.Name}
		cfg.Jira.Sprints = list
		cfg.Save()
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sprintUpdateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sprintUpdateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func askBoard() (*JiraBoard, error) {
	boards, err := getBoards()
	board := &JiraBoard{}
	if err != nil {
		printErr("error: %s", err)
		return nil, err
	}

	if len(boards) == 0 {
		return nil, fmt.Errorf("no boards found")
	}

	if len(boards) == 1 {
		board.ID = boards[0].ID
		board.Name = boards[0].Name
		return board, nil
	}

	boardsList := []string{}
	boardsMap := map[string]int{}
	for _, b := range boards {
		boardsList = append(boardsList, b.Name)
		boardsMap[b.Name] = b.ID
	}

	answer := ""
	q := &survey.Select{
		Message: "Choose a board:",
		Options: boardsList,
	}
	survey.AskOne(q, &answer, survey.Required)

	return &JiraBoard{ID: boardsMap[answer], Name: answer}, nil
}
