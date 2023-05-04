package geco_test

import (
	"os"
	"testing"

	"github.com/happyagosmith/geco/internal"
	"github.com/happyagosmith/geco/internal/geco"
	"gopkg.in/yaml.v2"
)

const errMsg = "got %d, want %d"

func TestGecoFile(t *testing.T) {
	t.Run("initialize", func(t *testing.T) {
		dir := t.TempDir()
		want := struct{ Templates []geco.Template }{Templates: []geco.Template{
			{
				URL:            "url1",
				Name:           "name1",
				GeneratedFiles: []string{"file1.1", "file1.2"},
				Model:          "model1",
			},
		}}
		d, _ := yaml.Marshal(want)
		err := os.WriteFile(dir+"/.geco.yml", d, 0600)
		if err != nil {
			t.Fatalf("err should be nil: %s", err.Error())
		}

		g, err := geco.NewGeco(geco.WithPath(dir))
		if err != nil {
			t.Fatalf("err should be nil: %s", err.Error())
		}
		got := g.GetTemplates()
		if len(got) != len(want.Templates) {
			t.Fatalf(errMsg, len(got), len(want.Templates))
		}
		assertTemplate(t, got[0], want.Templates[0])
	})

	t.Run("add templates", func(t *testing.T) {
		wants := []geco.Template{
			{
				URL:            "url1",
				Name:           "name1",
				GeneratedFiles: []string{"file1.1", "file1.2"},
				Model:          "model1",
			},
			{
				URL:            "url2",
				Name:           "name2",
				GeneratedFiles: []string{"file2.1", "file2.2"},
				Model:          "model2",
			},
		}
		g, _ := geco.NewGeco()
		g.AddTemplate(wants[0])
		g.AddTemplate(wants[1])
		got := g.GetTemplates()
		if len(got) != 2 {
			t.Fatalf(errMsg, len(got), 2)
		}
		for i, want := range wants {
			assertTemplate(t, got[i], want)
		}
	})

	t.Run("save file", func(t *testing.T) {
		dir := t.TempDir()
		want := geco.Template{
			URL:            "url1",
			Name:           "name1",
			GeneratedFiles: []string{"file1.1", "file1.2"},
			Model:          "model1",
		}
		g, _ := geco.NewGeco(geco.WithPath(dir))
		g.AddTemplate(want)

		g1, _ := geco.NewGeco(geco.WithPath(dir))
		got := g1.GetTemplates()
		if len(got) != 1 {
			t.Fatalf(errMsg, len(got), 1)
		}

		assertTemplate(t, got[0], want)
	})

	t.Run("overrides existing template", func(t *testing.T) {
		dir := t.TempDir()
		t1 := geco.Template{
			URL:            "url1",
			Name:           "name1",
			GeneratedFiles: []string{"file1.1", "file1.2"},
			Model:          "model1",
		}
		t2 := geco.Template{
			URL:            "url1",
			Name:           "name1",
			GeneratedFiles: []string{"file1.1"},
			Model:          "model1",
		}
		g, _ := geco.NewGeco(geco.WithPath(dir))
		g.AddTemplate(t1)
		g.AddTemplate(t2)

		got := g.GetTemplates()
		if len(got) != 1 {
			t.Fatalf(errMsg, len(got), 1)
		}

		assertTemplate(t, got[0], t2)
	})
}

func assertTemplate(t *testing.T, got, want geco.Template) {
	t.Helper()

	internal.AssertEqualString(t, "name", got.Name, want.Name)
	internal.AssertEqualString(t, "URL", got.URL, want.URL)
	internal.AssertEqualString(t, "Model", got.Model, want.Model)

	for j, wgf := range want.GeneratedFiles {
		internal.AssertEqualString(t, "generates file", got.GeneratedFiles[j], wgf)
	}
}
