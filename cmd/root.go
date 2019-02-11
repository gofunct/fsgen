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
	"github.com/gofunct/fsgen/pkg"
	"github.com/spf13/cobra"
	"os"
)

type Config struct {
	ConfigPath string
	OutputDir  string
	TmplDir    string
}

var V = pkg.NewViper()

var C *Config

func init() {
	{
		rootCmd.AddCommand(initCmd)
	}
	{
		rootCmd.PersistentFlags().StringVar(V.StringVarP(&C.ConfigPath, "p", ".", true, "FSGEN", "path to your config file:"))
		rootCmd.PersistentFlags().StringVar(V.StringVarP(&C.OutputDir, "o", "gen", true, "FSGEN", "path to target output directory:"))
		rootCmd.PersistentFlags().StringVar(V.StringVarP(&C.TmplDir, "t", "templates", true, "FSGEN", "path to your template directory:"))

	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fsgen",
	Short: "generate static assets for http.Filesystem",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
