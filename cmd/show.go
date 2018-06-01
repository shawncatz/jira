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
	"github.com/spf13/viper"
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

		fmt.Printf("%-15s: %+v\n", white(issue.Key), cyan(issue.Fields.Summary))
		fmt.Printf("%-15s: %s\n", white("Type"), cyan(issue.Fields.Type.Name))
		fmt.Printf("%-15s: %s\n", white("Priority"), cyan(issue.Fields.Priority.Name))
		if issue.Fields.Assignee != nil {
			fmt.Printf("%-15s: %s\n", white("Assigned"), cyan(issue.Fields.Assignee.Name))
		}
		fmt.Printf("%-15s:\n%s\n\n", white("Description"), issue.Fields.Description)
		if len(issue.Fields.Comments.Comments) > 0 {
			for _, c := range issue.Fields.Comments.Comments {
				fmt.Printf("%s (%s)\n%s\n\n", white(c.Author.Name), cyan(c.Author.EmailAddress), gray(c.Body))
			}
		}
		fmt.Println(cyan(viper.GetString("jira_base") + "/browse/" + issue.Key))
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
