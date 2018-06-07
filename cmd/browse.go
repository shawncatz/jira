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
	"os"
	"os/exec"
	"syscall"

	"github.com/spf13/cobra"
)

// browseCmd represents the browse command
var browseCmd = &cobra.Command{
	Use:   "browse <ISSUE>",
	Short: "Open a browser to JIRA issue",
	Long:  "Open a browser to JIRA issue",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		issue, _, err := jiraClient.Issue.Get(args[0], nil)
		if err != nil {
			printErr("error finding issue (%s):\n%s\n", args[0], err)
			return
		}
		url := issueURL(issue.Key)
		fmt.Println(cyan("opening: " + url))
		bin, err := exec.LookPath("open")
		if err != nil {
			printErr("error: %s", err.Error())
		}
		env := os.Environ()
		syscall.Exec(bin, []string{"open", url}, env)
	},
}

func init() {
	rootCmd.AddCommand(browseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// browseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// browseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
