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
	"github.com/jessevdk/go-assets"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"text/template"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "generate assets",
	Run: func(cmd *cobra.Command, args []string) {
		files := modules.Init
		for _, f := range files.Files {
			ProcessAsset(template.New("root"), f)
		}
	},
}

// Template reads a go template and writes it to dist given data.
func ProcessAsset(t *template.Template, file *assets.File) {
	if file.Name() == "/" {
		return
	}
	content := string(file.Data)

	tpl := t.New(file.Name())
	tpl, err := tpl.Parse(string(content))
	if err != nil {
		fatal("Could not parse template ", file.Name(), err)
	}

	f, err := os.Create(file.Name())
	if err != nil {
		fatal("Could not create file for writing ", file.Name(), err)
	}
	defer f.Close()
	err = tpl.Execute(f, viper.AllSettings())
	if err != nil {
		fatal("Could not execute template: ", file.Name(), err)
	}
}
