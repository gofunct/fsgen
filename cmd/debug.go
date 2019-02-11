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
	"github.com/spf13/viper"
	"io/ioutil"
)

func init() {
	rootCmd.AddCommand(debugCmd)
	debugCmd.Flags().BoolVar(&write, "write", false, "")
}

var (
	write bool
)

// debugCmd represents the debug command
var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		if write {
			pkg.L.FatalIfErr(ioutil.WriteFile("debug.txt", pkg.ToPrettyJson(viper.AllSettings()), 0755), cmd.UsageString(), "failed to write debug.txt file")
		} else {
			fmt.Println(pkg.ToPrettyJsonString(viper.AllSettings()))
		}
	},
}
