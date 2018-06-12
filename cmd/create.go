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

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a JIRA issue",
	Long:  "Create a JIRA issue",
	Run:   runCreate,
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runCreate(cmd *cobra.Command, args []string) {
	var questions = []*survey.Question{
		{
			Name:     "project",
			Prompt:   &survey.Input{Default: cfg.Jira.Project, Message: "Project for issue?"},
			Validate: survey.Required,
		},
		{
			Name: "type",
			Prompt: &survey.Select{
				Message: "Choose an issue type:",
				Options: cfg.Jira.Types,
				Default: cfg.DefaultType(),
			},
		},
		{
			Name: "sprint",
			Prompt: &survey.Select{
				Message: "Choose a sprint:",
				Options: cfg.Sprints(),
				Default: "Backlog",
			},
		},
		{
			Name:      "title",
			Prompt:    &survey.Input{Message: "Title for issue?"},
			Validate:  survey.Required,
			Transform: survey.Title,
		},
		{
			Name:   "description",
			Prompt: &survey.Editor{Message: "Please enter a description"},
		},
	}

	answers := &CreateAnswers{}
	err := survey.Ask(questions, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if debug {
		fmt.Printf("answers: %#v\n", answers)
	}

	issue, err := jiraCreate(answers)
	if err != nil {
		printErr("error: %s\n", err)
		return
	}

	fmt.Printf("Created: %s\n", white(issue.Key))
	fmt.Printf("%s\n", cyan(issueURL(issue.Key)))
}
