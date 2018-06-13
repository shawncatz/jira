package cmd

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfg = &Config{}

// Config stores the config
type Config struct {
	Jira JiraConfig
}

// JiraConfig stores the config
type JiraConfig struct {
	Base    string
	User    string `yaml:"user" viper:"user"`
	Pass    string `yaml:"pass" viper:"pass"`
	Project string
	Board   *JiraBoard
	Types   []string
	Sprints []*JiraSprint
}

// JiraSprint is a simple represenation of a sprint
type JiraSprint struct {
	ID   int
	Name string
}

// JiraBoard is a simple represenation of a sprint
type JiraBoard struct {
	ID   int
	Name string
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".jira" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".jira")
	}

	viper.AutomaticEnv() // read in environment variables that match

	err := viper.ReadInConfig()
	if err != nil {
		printErr("error reading config: %s\n", err)
	}

	// If a config file is found, read it in.
	if err := viper.Unmarshal(cfg); err != nil {
		fmt.Printf("Can't read config: (%#v) %s\n", cfg, err)
		os.Exit(1)
	}

	if len(cfg.Jira.Types) == 0 {
		printErr("you must include at least one Type in the configuration.\n" +
			"add a list of types to " + viper.ConfigFileUsed() + ".")
		os.Exit(1)
	}
}

// DefaultType returns the default type or empty string
func (c *Config) DefaultType() string {
	if len(c.Jira.Types) > 0 {
		return c.Jira.Types[0]
	}
	return ""
}

// Sprints returns the list of sprint names
func (c *Config) Sprints() []string {
	list := []string{}
	for _, s := range c.Jira.Sprints {
		list = append(list, s.Name)
	}
	return append([]string{c.DefaultSprint()}, list...)
}

// DefaultSprint returns the default sprint
func (c *Config) DefaultSprint() string {
	return "Backlog"
}

func (c *Config) findSprint(name string) *JiraSprint {
	for _, s := range c.Jira.Sprints {
		if s.Name == name {
			return s
		}
	}
	return nil
}

// Save the configuration to file
func (c *Config) Save() error {
	b, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(viper.ConfigFileUsed(), b, 0600)
	if err != nil {
		return err
	}

	return nil
}
