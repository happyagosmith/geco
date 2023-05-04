package geco

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Template struct {
	Name           string
	URL            string
	GeneratedFiles []string
	Model          string
}

type Geco struct {
	path      string
	Templates []Template
}

type Option func(f *Geco)

func WithPath(path string) Option {
	return func(g *Geco) {
		g.path = path
	}
}

const filename = ".geco.yml"

func NewGeco(opts ...Option) (Geco, error) {
	var g Geco
	for _, opt := range opts {
		opt(&g)
	}
	f, err := os.ReadFile(g.path + "/" + filename)
	if err != nil {
		return g, err
	}
	err = yaml.Unmarshal(f, &g)

	return g, err
}

func (g *Geco) AddTemplate(t Template) {
	for i, tmpl := range g.Templates {
		if tmpl.Name == t.Name {
			g.Templates[i].URL = t.URL
			g.Templates[i].GeneratedFiles = t.GeneratedFiles
			d, _ := yaml.Marshal(&g)
			_ = os.WriteFile(g.path+"/"+filename, d, 0600)

			return
		}
	}
	g.Templates = append(g.Templates, t)
	d, _ := yaml.Marshal(&g)
	_ = os.WriteFile(g.path+"/"+filename, d, 0600)
}

func (g Geco) GetTemplates() []Template {
	return g.Templates
}
