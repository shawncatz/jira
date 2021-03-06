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
	"sort"

	"github.com/spf13/cobra"
)

// sprintCmd represents the sprint command
var sprintCmd = &cobra.Command{
	Use:   "sprint",
	Short: "List sprints",
	Long:  "List sprints",
	Run: func(cmd *cobra.Command, args []string) {
		a, err := cmd.Flags().GetBool("all")
		if err != nil {
			printErr("failed to get option 'all'")
		}
		boards, err := getBoards()
		if err != nil {
			printErr("error finding boards: %s\n", err)
			return
		}
		for _, e := range boards {
			fmt.Printf("%d: (%s) %s\n", e.ID, e.Type, e.Name)
			sprints, err := getSprints(e.ID, a)
			if err != nil {
				printErr("error finding sprints: %s\n", err)
				return
			}

			sort.Sort(SprintCompletedDate(sprints))
			for _, s := range sprints {
				if a || s.State != "closed" {
					fmt.Printf("   %3d: (%s) %s\n", s.ID, s.EndDate, s.Name)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(sprintCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sprintCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	sprintCmd.Flags().BoolP("all", "a", false, "Show all sprints")
}
