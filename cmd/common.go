package cmd

import (
	"github.com/spf13/viper"
)

func issueURL(id string) string {
	return viper.GetString("jira.base") + "/browse/" + id
}
