// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
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
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
)

var (
	l, _  = zap.NewDevelopment()
	fatal = func(prefix string, msg string, err error) {
		l.Fatal(msg, zap.Namespace(prefix), zap.Error(err))
	}
	debug = func(prefix string, msg string, err error) {
		l.Debug(msg, zap.Namespace(prefix), zap.Error(err))
	}
	cfgFile string
	outDir  string
	walkDir string

	fs = afero.Afero{
		afero.NewOsFs(),
	}
)

func init() {
	zap.ReplaceGlobals(l)
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $PWD/.fsgen.yaml)")
	rootCmd.AddCommand(initCmd)
	rootCmd.PersistentFlags().StringVarP(&outDir, "output-dir", "o", "gen", "")
	rootCmd.PersistentFlags().StringVarP(&walkDir, "walk-dir", "w", "templates", "")
	ans, err := fs.Exists(outDir)
	debug("check if out directory exists", "main.init", err)
	if !ans {
		fatal("make all directories if they dont exist", "main.init", os.MkdirAll(outDir, 0755))
	}

}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fsgen",
	Short: "generate static assets for http.Filesystem",
	PersistentPostRun: func(cmd *cobra.Command, args []string) {

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {

		// Search config in home directory with name ".temp" (without extension).
		viper.AddConfigPath(os.Getenv("PWD"))
		viper.SetConfigName(".fsgen")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
