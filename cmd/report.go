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
	"time"
)

// reportCmd represents the report command
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Create a sprint report",
	Long:  "Create a sprint report",
	Run:   runReport,
}

func init() {
	rootCmd.AddCommand(reportCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type report struct {
	Done []*issue
	Todo []*issue
	Open []*issue
}

type issue struct {
	Key     string
	Points  float64
	Status  string
	Type    string
	Closed  time.Time
	Summary string
}

func runReport(cmd *cobra.Command, args []string) {
	a := ""
	p := &survey.Select{
		Message: "Select sprint: ",
		Options: cfg.Sprints(),
	}

	err := survey.AskOne(p, &a, survey.Required)
	if err != nil || a == "" {
		printErr("error getting answer: %s\n", err)
		return
	}

	report, err := getReport(a)
	if err != nil {
		printErr("%s\n", err)
		return
	}

	printReport(report)
}

func getReport(name string) (*report, error) {
	sprint := cfg.findSprint(name)
	if sprint == nil {
		return nil, fmt.Errorf("error finding sprint")
	}

	issues, err := getIssuesFromSprint(sprint.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting issues: %s", err)
	}

	field, err := getPointsField()
	if err != nil {
		return nil, fmt.Errorf("error getting points field: %s", err)
	}

	report := &report{}
	for _, i := range issues {
		if i.Fields.Type.Name == "Sub-task" {
			continue
		}

		var points float64
		if i.Fields.Unknowns[field.Key] != nil {
			points = i.Fields.Unknowns[field.Key].(float64)
		}

		issue := &issue{
			Key:     i.Key,
			Points:  points,
			Status:  i.Fields.Status.Name,
			Type:    i.Fields.Type.Name,
			Closed:  time.Time(i.Fields.Resolutiondate),
			Summary: i.Fields.Summary,
		}

		switch i.Fields.Status.StatusCategory.Name {
		case "Done":
			report.Done = append(report.Done, issue)
		case "To Do":
			report.Todo = append(report.Todo, issue)
		case "In Progress":
			report.Open = append(report.Open, issue)
		}

	}

	return report, nil
}

func printReport(report *report) {
	fmt.Printf("\n%s\n", white("To Do"))
	for _, i := range report.Todo {
		printIssue(i)
	}
	fmt.Printf("\n%s\n", white("Open"))
	for _, i := range report.Open {
		printIssue(i)
	}
	fmt.Printf("\n%s\n", white("Done"))
	for _, i := range report.Done {
		printIssue(i)
	}
}

func printIssue(issue *issue) {
	fmt.Printf("%10.10s %3.0f %-15.15s %-10.10s %-50.50s\n", cyan(issue.Key), issue.Points, issue.Status, issue.Type, issue.Summary)
}
