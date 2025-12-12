/*
Copyright © 2025 Artur Taranchiev <artur.taranchiev@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "ffiii-rate-updater",
	Short: "A tool to update exchange rates in Firefly III",
	Long: `ffiii-rate-updater is a command-line tool that fetches exchange rates for specified currencies
and updates them in Firefly III via its API. For example:

    ffiii-rate-updater update

or you can set currencies and date via flags:

    ffiii-rate-updater update --currencies USD,EUR,GBP --date 2025-01-01`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initializeConfig(cmd)
	},
}

var initConfigCmd = &cobra.Command{
	Use:   "init-config",
	Short: "Generate a default configuration file",
	Long:  `Generate a default configuration file for ffiii-rate-updater.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		initViper := viper.New()
		initViper.Set("firefly.api_key", viper.GetString("firefly.api_key"))
		initViper.Set("firefly.api_url", viper.GetString("firefly.api_url"))
		initViper.Set("currencies", viper.GetStringSlice("currencies"))

		initViper.AddConfigPath(".")
		initViper.SetConfigName("config")
		initViper.SetConfigType("yaml")
		initViper.SetConfigFile("./config.yaml")

		err := initViper.SafeWriteConfig()
		if err != nil {
			var configFileAlreadyExistsError viper.ConfigFileAlreadyExistsError
			if errors.As(err, &configFileAlreadyExistsError) {
				return err
			}
		}

		fmt.Println("Configuration file created at:", initViper.ConfigFileUsed())
		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/ffiii-rate-updater/config)")
	rootCmd.PersistentFlags().StringP("firefly.api_key", "k", "your_firefly_api_key_here", "Firefly III API key")
	rootCmd.PersistentFlags().StringP("firefly.api_url", "u", "https://your-firefly-iii-instance.com/api/v1", "Firefly III API URL")
	rootCmd.PersistentFlags().StringSliceP("currencies", "c", []string{}, "List of currencies to fetch exchange rates for (e.g. USD,EUR,GBP)")
	rootCmd.PersistentFlags().StringP("date", "d", "latest", "Date for which to fetch exchange rates (format: YYYY-MM-DD or 'latest')")

	rootCmd.AddCommand(initConfigCmd)
}

func initializeConfig(cmd *cobra.Command) error {

	viper.SetEnvPrefix("FFIII_RATE_UPDATER")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "*", "-", "*"))
	viper.AutomaticEnv()

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(".")
		viper.AddConfigPath(home + "/.config/ffiii-rate-updater")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return err
		}
	} else {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	err := viper.BindPFlags(cmd.Flags())
	if err != nil {
		return err
	}

	return nil
}
