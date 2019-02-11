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
// limitations under the License. //go:generate go-assets-builder -p modules -s="/init" -o init.go -v Init init

package cmd

import (
	"github.com/gofunct/fsgen/modules"
	"github.com/gofunct/fsgen/pkg"
	"github.com/spf13/cobra"
	"text/template"
)

type project struct {
	Contributors []string
	GHUserName   string
	Description  string
	DHUserName   string
	BaseImage    string
	BuildImage   string
}

var Project *project

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "generate initial project assets",
	Run: func(cmd *cobra.Command, args []string) {
		files := modules.Init
		pkg.L.WarnIfErr(V.ReadInConfig(), V.ConfigFileUsed(), "failed to read in config")
		for _, f := range files.Files {
			V.ProcessAsset(template.New("root"), f)
		}
	},
}
