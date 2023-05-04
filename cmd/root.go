/*
Copyright © 2023 Happy Smith happyagosmith@gmail.com

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
	"fmt"
	"os"

	"github.com/happyagosmith/geco/internal/geco"
	"github.com/happyagosmith/geco/internal/renderer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type GecoFile interface {
	AddTemplate(t geco.Template)
	GetTemplates() []geco.Template
}

var cfgFile string

func NewRootCmd(fg renderer.FilesGenerator, gf GecoFile) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "geco",
		Short: "A brief description of your application",
		Long:  `A longer description that spans multiple lines`,
	}

	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.geco.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.AddCommand(newApplyCmd(fg, gf))
	rootCmd.AddCommand(newInitCmd(fg, gf))
	rootCmd.AddCommand(newUpdateCmd(fg, gf))

	return rootCmd
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".geco")
	}

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
