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
	"os"

	"github.com/spf13/cobra"
)

var defaultConfig = []byte(`
jira:
  base: https://yourcompany.atlassian.net
  user: user@email.com       
  pass: password_or_api_key  # can be one or the other, api key strongly recommended
  project: PROJECT           # your default project
  images: false              # does your terminal support images?
  types:                     # hope to automate management of these soon
    - Bug
    - Task
    - Story
  sprints:
  board:
`)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Generate a configuration file",
	Long:  "Generate a configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		path := os.Getenv("HOME") + "/.jira.yaml"

		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			printErr("Error: getting flag: %s", err.Error())
			return
		}

		if _, err = os.Stat(path); err == nil && !force {
			printErr("Error: %s already exists", path)
			return
		}

		err = ioutil.WriteFile(path, defaultConfig, 0600)
		if err != nil {
			printErr("Error: %s\n", err.Error())
			return
		}

		fmt.Printf("file %s created.\nAdd your jira configuration, then run %s to update the config from Jira.\n", white(path), white("jira update"))
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	initCmd.Flags().BoolP("force", "f", false, "force (overwrite existing file)")
}
