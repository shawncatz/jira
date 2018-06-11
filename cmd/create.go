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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	defaultProject := viper.GetString("jira.project")

	typeOptions := viper.GetStringSlice("jira.types")
	typeDefault := ""
	if len(typeOptions) > 0 {
		typeDefault = typeOptions[0]
	} else {
		printErr("you must include at least one Type in the configuration.\n" +
			"add a list of types to " + viper.ConfigFileUsed() + ".")
		return
	}

	sprintOptions := []string{"Backlog"}
	sprintOptions = append(sprintOptions, viper.GetStringSlice("jira.sprints")...)

	sprintDefault := "Backlog"
	if len(sprintOptions) > 0 {
		sprintDefault = sprintOptions[0]
	}

	var questions = []*survey.Question{
		{
			Name:     "project",
			Prompt:   &survey.Input{Default: defaultProject, Message: "Project for issue?"},
			Validate: survey.Required,
		},
		{
			Name: "type",
			Prompt: &survey.Select{
				Message: "Choose an issue type:",
				Options: typeOptions,
				Default: typeDefault,
			},
		},
		{
			Name: "sprint",
			Prompt: &survey.Select{
				Message: "Choose a sprint:",
				Options: sprintOptions,
				Default: sprintDefault,
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
	answers := struct {
		Project     string
		Title       string
		Description string
		Type        string
		Sprint      string
	}{}

	err := survey.Ask(questions, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if debug {
		fmt.Printf("answers: %#v\n", answers)
	}

	f := &jira.IssueFields{
		Project:     jira.Project{Key: answers.Project},
		Type:        jira.IssueType{Name: answers.Type},
		Summary:     answers.Title,
		Description: answers.Description,
		Labels:      []string{"from-cli"},
	}

	if answers.Sprint != "Backlog" {
		f.Sprint = &jira.Sprint{Name: answers.Sprint}
	}

	i := jira.Issue{
		Fields: f,
	}

	if debug {
		fmt.Printf("%#v\n", i)
		fmt.Printf("%#v\n", i.Fields)
	}

	issue, response, err := jiraClient.Issue.Create(&i)
	if err != nil {
		printErr(err.Error())
		b, _ := ioutil.ReadAll(response.Response.Body)
		fmt.Printf("response:\n%s\n", string(b))
		return
	}

	fmt.Printf("Created: %s\n", white(issue.Key))
	fmt.Printf("%s\n", cyan(issueURL(issue.Key)))
}
