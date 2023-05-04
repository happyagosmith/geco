package renderer

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/happyagosmith/geco/internal/format"
	"gopkg.in/yaml.v3"
)

type Options struct {
	OPath string
	MPath string
	TPath string
}

type FilesGenerator interface {
	Run(opts Options, out io.Writer) ([]string, error)
}

type Renderer struct{}

func (r Renderer) Run(opts Options, out io.Writer) ([]string, error) {
	tInfo, err := os.Stat(opts.TPath)
	if os.IsNotExist(err) {
		return []string{}, fmt.Errorf("template %s not exists", opts.TPath)
	}

	if _, err := os.Stat(opts.MPath); os.IsNotExist(err) {
		return []string{}, fmt.Errorf("model %s not exists", opts.MPath)
	}

	dir := filepath.Dir(opts.OPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return []string{}, fmt.Errorf("output %s not exists", dir)
	}

	if notDir := !tInfo.IsDir(); notDir {
		if opts.OPath != "" {
			name := filepath.Base(opts.TPath)
			name = strings.Replace(name, ".tpl", "", 1)
			opts.OPath = fmt.Sprintf("%s/%s", opts.OPath, name)
		}
		gFile, err := executeOnFile(opts, out)
		if err != nil {
			return []string{}, err
		}

		return []string{gFile}, nil
	}

	return executeOnFolder(opts, out)
}

func executeOnFile(o Options, out io.Writer) (string, error) {
	tpl, err := readTemplate(o)
	if err != nil {
		return "", err
	}

	model, err := readModel(o.MPath)
	if err != nil {
		return "", err
	}

	if o.OPath == "" {
		return "", tpl.Execute(out, model)
	}

	err = os.MkdirAll(filepath.Dir(o.OPath), 0700)
	if err != nil {
		return "", err
	}
	f, err := os.Create(o.OPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	return o.OPath, tpl.Execute(f, model)
}

func executeOnFolder(o Options, out io.Writer) ([]string, error) {
	var gFiles []string
	name := filepath.Base(o.TPath)
	err := filepath.Walk(o.TPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.Contains(info.Name(), ".tpl") {
			return nil
		}

		rPath := strings.Join(strings.Split(path, name)[1:], name)
		rPath = strings.Replace(rPath, ".tpl", "", 1)
		gFile, err := executeOnFile(Options{TPath: path, MPath: o.MPath, OPath: o.OPath + rPath}, out)
		if err != nil {
			return err
		}
		gFiles = append(gFiles, gFile)

		return nil
	})

	return gFiles, err
}

func readTemplate(o Options) (*template.Template, error) {
	funcMap := template.FuncMap{
		"toYaml": format.ToYaml,
		"mergeYaml": func(a, b *format.Yaml) (string, error) {
			if b == nil {
				return a.String()
			}

			if err := a.Merge(*b); err != nil {
				return "", err
			}
			return a.String()
		},
		"readYaml": func(path string) (*format.Yaml, error) {
			filename := filepath.Dir(o.OPath) + "/" + path
			_, err := os.Stat(filename)
			if os.IsNotExist(err) {
				return nil, nil
			}

			rf, err := ioutil.ReadFile(filename)
			if err != nil {
				return nil, err
			}

			y, err := format.NewYaml(rf)
			if err != nil {
				return nil, err
			}

			return &y, err
		},
		"readTemplateYaml": func(path string) (*format.Yaml, error) {
			filename := filepath.Dir(o.TPath) + "/" + path
			_, err := os.Stat(filename)
			if os.IsNotExist(err) {
				return nil, nil
			}

			o := Options{TPath: filename, MPath: o.MPath, OPath: ""}
			var b bytes.Buffer
			_, err = executeOnFile(o, &b)
			if err != nil {
				return nil, err
			}
			y, err := format.NewYaml(b.Bytes())
			if err != nil {
				return nil, err
			}

			return &y, err
		},
		"extractSections": func(y *format.Yaml) ([]format.Section, error) {
			if y == nil {
				return []format.Section{}, nil
			}
			return y.ParseComments()
		},
	}
	tpl, err := template.New(filepath.Base(o.TPath)).Delims("[[", "]]").Funcs(sprig.FuncMap()).Funcs(funcMap).ParseFiles(o.TPath)

	return tpl, err
}

func readModel(vPath string) (map[string]any, error) {
	v, err := os.ReadFile(vPath)
	if err != nil {
		return nil, err
	}
	var values map[string]any
	err = yaml.Unmarshal(v, &values)

	return values, err
}
