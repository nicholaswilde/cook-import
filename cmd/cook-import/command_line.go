package main

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var version string

func possibleLogLevels() []string {
	levels := make([]string, 0)

	for _, l := range log.AllLevels {
		levels = append(levels, l.String())
	}

	return levels
}

func initializeCli() {
	logLevelName := viper.GetString("log-level")
	logLevel, err := log.ParseLevel(logLevelName)
	if err != nil {
		log.Errorf("Failed to parse provided log level %s: %s", logLevelName, err)
		os.Exit(1)
	}

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetLevel(logLevel)
}

func newCookImportCommand(run func(cmd *cobra.Command, args []string)) (*cobra.Command, error) {
	command := &cobra.Command{
		Use:     "cook-import",
		Short:   "cook-import is a command line tool to import recipes into Cooklang format using ChatGPT",
		Version: version,
		Run:     run,
	}
	logLevelUsage := fmt.Sprintf("level of logs that should printed, one of (%s)", strings.Join(possibleLogLevels(), ", "))
	command.PersistentFlags().BoolP("file", "f", false, "If you want the output to be in a file, use this flag. Otherwise defaults to console screen.")
	command.PersistentFlags().StringP("log-level", "L", "info", logLevelUsage)
	command.PersistentFlags().StringP("openai-api-key", "k", "", "OpenAI API key")
	command.PersistentFlags().StringP("link", "l", "", "Input a url link to a recipe")
	command.MarkPersistentFlagRequired("link")
	
	viper.SetConfigName(".cookimport")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/cook-import/")
	viper.AddConfigPath("$HOME/.config/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Debug(err)
		} else {
			log.Debugf("Error occurred while reading config file: %s \n", err)
		}
	} else {
		log.Debugf("Using config file %s", viper.ConfigFileUsed())
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("COOK_IMPORT")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.RegisterAlias("dryRun", "dry-run")
	viper.RegisterAlias("logLevel", "log-level")
	viper.RegisterAlias("openAiApiKey", "openai-api-key")
	err = viper.BindPFlags(command.PersistentFlags())
	return command, err
}
