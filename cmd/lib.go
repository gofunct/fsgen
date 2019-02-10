package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/jessevdk/go-assets"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Template reads a go template and writes it to dist given data.
func ProcessAsset(t *template.Template, file *assets.File) {
	if file.Name() == "/" {
		return
	}
	content := string(file.Data)

	tpl := t.New(file.Name()).Funcs(sprig.GenericFuncMap())
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

// toPrettyJson encodes an item into a pretty (indented) JSON string
func ToPrettyJsonString(v interface{}) string {
	output, _ := json.MarshalIndent(v, "", "  ")
	return string(output)
}

// toPrettyJson encodes an item into a pretty (indented) JSON string
func ToPrettyJson(v interface{}) []byte {
	output, _ := json.MarshalIndent(v, "", "  ")
	return output
}

// Prompt prompts user for input with default value.
func Prompt(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	return text
}

func WalkTemplates(dir string) {
	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			debug(path, "walkfunc copy error", err)
		}
		if strings.Contains(path, ".tmpl") {
			b, err := ioutil.ReadFile(path)
			newt, err := template.New(info.Name()).Funcs(sprig.GenericFuncMap()).Parse(string(b))
			if err != nil {
				return err
			}

			f, err := os.Create(outDir + strings.TrimSuffix(info.Name(), ".tmpl"))
			if err != nil {
				return err
			}
			return newt.Execute(f, viper.AllSettings())
		}
		return nil
	}); err != nil {
		fatal(dir, "failed to walk directory", err)
	}
}

func CopyFile(srcfile, dstfile string) (*os.File, error) {
	srcF, err := os.Open(srcfile) // nolint: gosec
	if err != nil {
		return nil, fmt.Errorf("could not open source file: %s", err)
	}
	defer srcF.Close()

	dstF, err := os.Create(dstfile)
	if err != nil {
		return nil, err
	}

	if _, err = io.Copy(dstF, srcF); err != nil {
		return nil, fmt.Errorf("could not copy file: %s", err)
	}
	return dstF, os.Chmod(dstfile, 0755)
}
