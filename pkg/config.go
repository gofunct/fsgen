package pkg

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/jessevdk/go-assets"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tcnksm/go-input"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	Fs = afero.Afero{
		afero.NewOsFs(),
	}
)

type Viper struct {
	Q *input.UI
	*viper.Viper
}

func NewViper() *Viper {
	v := viper.GetViper()
	v.SetFs(Fs)
	return &Viper{Viper:v, Q: input.DefaultUI()}
}

func (v *Viper) StringVarP(ptr *string, key string, def string, required bool, envprefix string, question string) (*string, string, string, string) {
	var (
		val string
	)
	if envprefix != "" {
		v.SetEnvPrefix(envprefix)
	}
	if def != "" {
		v.SetDefault(key, def)
	}

	if v.GetString(key) == "" && required {
		_, err := v.Q.Select(question, []string{"t", "f"}, &input.Options{
			Default:  def,
			Loop:     required,
			Required: required,
			ValidateFunc: func(s string) error {
				if len(s) > 50 {
					return errors.New("input must be less than 50 characters expected")
				}
				if s == "" && required {
					return errors.New("input must not be empty")
				}
				val = s
				v.Set(key, val)
				_ = os.Setenv(envprefix+"_"+strings.ToUpper(key), val)
				return nil
			},
		})
		if err != nil {
			L.DebugIfErr(err, question, "failed to query selection")
			return ptr, key, val, question
		}

	}

	if v.InConfig(key) {
		val = v.GetString(key)
	}

	return ptr, key, val, question
}
func (v *Viper) BoolVarP(ptr *bool, key string, def bool, required bool, envprefix string, question string) (*bool, string, bool, string) {
	var (
		val bool
	)
	if envprefix != "" {
		v.SetEnvPrefix(envprefix)
	}
	if def != false {
		v.SetDefault(key, def)
	}

	var defstring string
	if def {
		defstring = "t"
	} else {
		defstring = "f"
	}
	if !v.InConfig(key) {
		_, err := v.Q.Select(question, []string{"t", "f"}, &input.Options{
			Default:  defstring,
			Loop:     required,
			Required: required,
			ValidateFunc: func(s string) error {
				if s == "t" {
					val = true
					return nil
				}
				if s == "f" {
					val = false
					return nil
				}
				return errors.New("must provide either t(true) or f(false) for bool variable")
			},
		})
		if err != nil {
			L.DebugIfErr(err, question, "failed to query selection")
			return ptr, key, val, question
		}

	}

	if v.InConfig(key) {
		val = v.GetBool(key)
	}

	return ptr, key, val, question
}

func (v *Viper) Init(root *cobra.Command, cfgFile string) func() {
	return func() {
		if v == nil {
			v = &Viper{
				Viper: viper.New(),
			}
			v.SetFs(Fs)
			v.AutomaticEnv() // read in environment variables that match
		}
		if cfgFile != "" {
			// Use config file from the flag.
			v.SetConfigFile(cfgFile)
		} else {

			// Search config in home directory with name ".temp" (without extension).
			v.AddConfigPath(os.Getenv("PWD"))
			v.SetConfigName(".config")
			v.SetConfigType("yaml")
		}

		// If a config file is found, read it in.
		if err := v.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", v.ConfigFileUsed())
		}

		if root.HasAvailableFlags() {
			L.DebugIfErr(v.BindPFlags(root.Flags()), root.Name(), "failed to bind config to Flags()")
		}
		if root.HasAvailablePersistentFlags() {
			L.DebugIfErr(v.BindPFlags(root.PersistentFlags()), root.Name(), "failed to bind config to PersistentFlags()")
		}
		if root.HasAvailableLocalFlags() {
			L.DebugIfErr(v.BindPFlags(root.LocalFlags()), root.Name(), "failed to bind config to LocalFlags()")
			L.DebugIfErr(v.BindPFlags(root.LocalNonPersistentFlags()), root.Name(), "failed to bind config to LocalNonPersistentFlags()")
		}
		if root.HasAvailableInheritedFlags() {
			L.DebugIfErr(v.BindPFlags(root.InheritedFlags()), root.Name(), "failed to bind config to InheritedFlags()")
		}
		if root.HasAvailableSubCommands() {
			for _, cmd := range root.Commands() {

				if cmd.HasAvailableFlags() {
					L.DebugIfErr(v.BindPFlags(cmd.Flags()), cmd.Name(), "failed to bind config to Flags()")
				}
				if cmd.HasAvailablePersistentFlags() {
					L.DebugIfErr(v.BindPFlags(cmd.PersistentFlags()), cmd.Name(), "failed to bind config to PersistentFlags()")
				}
				if cmd.HasAvailableLocalFlags() {
					L.DebugIfErr(v.BindPFlags(cmd.LocalFlags()), cmd.Name(), "failed to bind config to LocalFlags()")
					L.DebugIfErr(v.BindPFlags(cmd.LocalNonPersistentFlags()), cmd.Name(), "failed to bind config to LocalNonPersistentFlags()")
				}
				if cmd.HasAvailableInheritedFlags() {
					L.DebugIfErr(v.BindPFlags(cmd.InheritedFlags()), cmd.Name(), "failed to bind config to InheritedFlags()")
				}

			}
		}
		v.Sync()
	}
}

func (v *Viper) Write(p []byte) (n int, err error) {
	if v.ConfigFileUsed() == "" {
		bits, err := yaml.Marshal(v.AllSettings())
		L.FatalIfErr(err, "config", "failed to marshal config settings")
		L.FatalIfErr(ioutil.WriteFile(".config.yaml", bits, 0755), "config", "failed to write config settings to file")
		v.SetConfigName(".config")
		v.AddConfigPath(os.Getenv("PWD"))
		v.SetConfigType("yaml")
	}

	f, err := Fs.Open(v.ConfigFileUsed())
	L.FatalIfErr(err, v.ConfigFileUsed(), "failed to open config file")
	bits, err := yaml.Marshal(p)
	return io.WriteString(f, fmt.Sprintln(bits))
}

func (v *Viper) Read(p []byte) (int, error) {
	bits, _ := yaml.Marshal(p)
	s := strings.NewReader(string(bits))
	err := v.ReadConfig(s)
	return len(p), err
}

func (v *Viper) Sync() {
	for _, e := range os.Environ() {
		sp := strings.Split(e, "=")
		v.SetDefault(strings.ToLower(sp[0]), sp[1])
	}
	for k, v := range v.AllSettings() {
		val, ok := v.(string)
		if ok {
			L.DebugIfErr(os.Setenv(strings.ToUpper(k), val), k, "failed to bind "+val)
		}
	}
}

func (v *Viper) JsonSettings() []byte {
	return (ToPrettyJson(v.AllSettings()))
}

func (v *Viper) JsonSettingsString() string {
	return (ToPrettyJsonString(v.AllSettings()))
}

func (v *Viper) YamlSettings() []byte {
	bits, err := yaml.Marshal(v.AllSettings())
	L.WarnIfErr(err, v.ConfigFileUsed(), "failed to unmarshal config to yaml")
	return bits
}

// Prompt prompts user for input with default value.
func (v *Viper) Prompt(key, question string) string {
	switch {
	case v.InConfig(key):
		return v.GetString(key)
	case os.Getenv(strings.ToUpper(key)) != "":
		return os.Getenv(strings.ToUpper(key))
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print(question)
	text, _ := reader.ReadString('\n')
	v.Set(key, text)
	return text
}

// Prompt prompts user for input with default value.
func (v *Viper) PromptSet(key string, question string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(question)
	text, _ := reader.ReadString('\n')
	_ = os.Setenv(strings.ToUpper(key), text)
	v.Set(key, text)
}

// Prompt prompts user for input with default value.
func (v *Viper) PromptCSV(key string, question string) []string {
	switch {
	case v.InConfig(key):
		return v.GetStringSlice(key)
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(question)
	text, _ := reader.ReadString('\n')
	txtCsv, err := AsCSV(text)
	L.DebugIfErr(err, "prompt csv", "failed to read comma seperated values from input")
	v.Set(key, txtCsv)
	return txtCsv
}

// Prompt prompts user for input with default value.
func (v *Viper) PromptSetCSV(key string, question string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(question)
	text, _ := reader.ReadString('\n')
	txtCsv, err := AsCSV(text)
	L.DebugIfErr(err, "prompt csv", "failed to read comma seperated values from input")
	v.Set(key, txtCsv)
}

// Prompt prompts user for input with default value.
func (v *Viper) PromptMap(key string, question string) map[string]string {
	switch {
	case v.InConfig(key):
		return v.GetStringMapString(key)
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(question)
	text, _ := reader.ReadString('\n')
	txtMap, err := AsMap(text)
	L.DebugIfErr(err, "prompt map", "failed to read comma seperated values from input, seperate map values with : or = and map entries with ,")
	return txtMap
}

// Prompt prompts user for input with default value.
func (v *Viper) PromptSetMap(key string, question string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(question)
	text, _ := reader.ReadString('\n')
	txtMap, err := AsMap(text)
	L.DebugIfErr(err, "prompt map", "failed to read comma seperated values from input, seperate map values with : or = and map entries with ,")
	v.Set(key, txtMap)
}


// Template reads a go template and writes it to dist given data.
func (v *Viper) ProcessAsset(t *template.Template, file *assets.File) {
	if file.Name() == "/" {
		return
	}
	content := string(file.Data)

	tpl := t.New(file.Name()).Funcs(sprig.GenericFuncMap())
	tpl, err := tpl.Parse(string(content))
	if err != nil {
		L.WarnIfErr(err, file.Name(), "Could not parse template ")
	}

	f, err := Fs.Create(file.Name())
	if err != nil {
		L.WarnIfErr(err, file.Name(), "Could not create file for writing")
	}
	defer f.Close()
	err = tpl.Execute(f, v.AllSettings())
	if err != nil {
		L.WarnIfErr(err, file.Name(), "Could not execute template")
	}
}

func (v *Viper) WalkTemplates(dir string, outDir string) {

	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			L.DebugIfErr(err, path, "walkfunc copy error")
		}
		if strings.Contains(path, ".tmpl") {
			b, err := ioutil.ReadFile(path)
			newt, err := template.New(info.Name()).Funcs(sprig.GenericFuncMap()).Parse(string(b))
			if err != nil {
				return err
			}

			f, err := Fs.Create(outDir +"/" +strings.TrimSuffix(info.Name(), ".tmpl"))
			if err != nil {
				return err
			}
			return newt.Execute(f, v.AllSettings())
		}
		return nil
	}); err != nil {
		L.WarnIfErr(err, dir+" to "+outDir, "failed to walk templates")
	}
}

func (v *Viper) CopyFile(srcfile, dstfile string) (*afero.File, error) {
	srcF, err := Fs.Open(srcfile) // nolint: gosec
	if err != nil {
		return nil, fmt.Errorf("could not open source file: %s", err)
	}
	defer srcF.Close()

	dstF, err := Fs.Create(dstfile)
	if err != nil {
		return nil, err
	}

	if _, err = io.Copy(dstF, srcF); err != nil {
		return nil, fmt.Errorf("could not copy file: %s", err)
	}
	return &dstF, Fs.Chmod(dstfile, 0755)
}
