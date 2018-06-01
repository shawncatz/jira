// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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

	"github.com/andygrunwald/go-jira"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/AlecAivazis/survey.v1"
)

var questions = []*survey.Question{
	{
		Name:     "project",
		Prompt:   &survey.Input{Default: "FOND", Message: "Project for issue?"},
		Validate: survey.Required,
	},
	{
		Name: "type",
		Prompt: &survey.Select{
			Message: "Choose an issue type:",
			Options: []string{"Story", "Task", "Bug"},
			Default: "Story",
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
	//{
	//	Name: "Sprint",
	//	Prompt: &survey.Select{
	//		Message: "Choose a sprint:",
	//		Options: []string{"Backlog", "Candidates"},
	//		Default: "Backlog",
	//	},
	//},
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a JIRA issue",
	Long:  "Create a JIRA issue",
	Run: func(cmd *cobra.Command, args []string) {
		answers := struct {
			Project     string
			Title       string
			Description string
			Type        string
		}{}
		err := survey.Ask(questions, &answers)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Printf("answers: %#v\n", answers)

		i := jira.Issue{
			Fields: &jira.IssueFields{
				Reporter:    &jira.User{Name: viper.GetString("jira_account")},
				Project:     jira.Project{Key: answers.Project},
				Summary:     answers.Title,
				Description: answers.Description,
				Labels:      []string{"from-cli"},
			},
		}

		fmt.Printf("%#v\n", i)
		fmt.Printf("%#v\n", i.Fields)

		issue, _, err := jiraClient.Issue.Create(&i)
		if err != nil {
			fmt.Print(err.Error())
			return
		}

		fmt.Printf("Created: %s\n", white(issue.Key))
		fmt.Printf("%s", cyan(viper.GetString("jira_base")+"/browse/"+issue.Key))
	},
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
