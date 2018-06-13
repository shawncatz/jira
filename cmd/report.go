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
	"github.com/shawncatz/jira/report"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
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

	sprint := cfg.findSprint(a)
	if sprint == nil {
		printErr("error finding sprint")
		return
	}

	r := report.NewSprintReport(jiraClient, sprint.ID, sprint.Name)
	if err != nil {
		printErr("%s\n", err)
		return
	}

	err = r.Build()
	if err != nil {
		printErr("error building report: %s\n", err)
		return
	}

	r.Print()
}
