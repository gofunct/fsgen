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
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// compileCmd represents the compile command
var compileCmd = &cobra.Command{
	Use:   "compile",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		if err := filepath.Walk("templates", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				debug(path, "walkfunc copy error", err)
			}
			if strings.Contains(path, ".tmpl") {
				f, err := CopyFile(path, outDir+"/"+strings.TrimSuffix(filepath.Base(path), ".tmpl"))
				if err != nil {
					return err
				}
				b, err := ioutil.ReadFile(path)
				newt, err := template.New(info.Name()).Parse(string(b))
				if err != nil {
					return err
				}

				return newt.Execute(f, viper.AllSettings())
			}
			return nil
		}); err != nil {
			fatal(cmd.Name(), "failed to walk directory", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(compileCmd)

}

func CopyFile(srcfile, dstfile string) (*os.File, error) {
	srcF, err := os.Open(srcfile) // nolint: gosec
	if err != nil {
		return nil, fmt.Errorf("could not open source file: %s", err)
	}
	defer srcF.Close()
	 ans, _ := afero.Exists(fs, dstfile)

	 if ans {
	 	if err :=  os.Remove(dstfile); err != nil {
	 		return nil, err
		}
	 }
	dstF, err := os.Create(dstfile)
	if err != nil {
		return nil, err
	}


	if _, err = io.Copy(dstF, srcF); err != nil {
		return nil, fmt.Errorf("could not copy file: %s", err)
	}
	return dstF, os.Chmod(dstfile, 0755)
}
