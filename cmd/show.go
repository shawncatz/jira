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

	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Display information about a JIRA issue",
	Long:  "Display information about a JIRA issue",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		issue, _, err := jiraClient.Issue.Get(args[0], nil)
		if err != nil {
			fmt.Printf("error finding issue (%s):\n%s\n", args[0], err)
			return
		}

		fmt.Printf("%s: %+v\n", issue.Key, issue.Fields.Summary)
		fmt.Printf("Type: %s\n", issue.Fields.Type.Name)
		fmt.Printf("Priority: %s\n", issue.Fields.Priority.Name)
		if issue.Fields.Assignee != nil {
			fmt.Printf("Assigned: %s\n", issue.Fields.Assignee.Name)
		}
		fmt.Printf("Description:\n%s\n", issue.Fields.Description)
		if len(issue.Fields.Comments.Comments) > 0 {
			for _, c := range issue.Fields.Comments.Comments {
				fmt.Printf("--- %s\n%s\n", c.Author.Name, c.Body)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(showCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
