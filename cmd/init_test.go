package cmd_test

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/happyagosmith/geco/cmd"
	"github.com/happyagosmith/geco/internal"
	"github.com/happyagosmith/geco/internal/geco"
	"github.com/happyagosmith/geco/internal/renderer"
	shell "github.com/mattn/go-shellwords"
)

func TestInit(t *testing.T) {
	dest := t.TempDir()
	tests := []testCase{
		{
			name: "folder template",
			cmd:  "init testdata/template1 " + dest,
			wantOpts: renderer.Options{
				MPath: dest + "/geco-model.yaml",
				TPath: "testdata/template1",
				OPath: dest,
			},
		},
	}

	runCheckFilesTests(t, tests)
}

func TestModelFile(t *testing.T) {
	path := t.TempDir()
	gecoModel, _ := os.ReadFile("testdata/geco-model.yaml")
	_ = os.WriteFile(path + "/geco-model.yaml", gecoModel, 0644)

	c := fmt.Sprintf("init testdata/template1 %s", path)
	args, _ := shell.Parse(c)
	mg := mockGeco{addTemplate: func(gotTemplate geco.Template) { return }}
	fg := mockFileGenerator{run: func(got renderer.Options, out io.Writer) ([]string, error) {
		return nil, nil
	}}
	rootCmd := cmd.NewRootCmd(&fg, &mg)
	rootCmd.SetArgs(args)

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Got %s, expected nil", err.Error())
	}
	want, _ := os.ReadFile("testdata/want-geco-model.yaml")
	got, _ := os.ReadFile(path + "/geco-model.yaml")
	
	internal.AssertEqualString(t, "content of geco-model.yaml", string(got), string(want))
}

func TestGecoFile(t *testing.T) {
	wantTemplate := geco.Template{
		Name: "template1",
		URL:  "testdata/template1",
		GeneratedFiles: []string{
			"./template.md",
			"./nestedFolder/template.md",
			"./nestedFolder/template1.md",
		},
		Model: "./geco-model.yaml",
	}

	tests := []struct {
		name     string
		cmd      string
		gecoFile mockGeco
	}{
		{
			name: "first template",
			cmd:  "init testdata/template1 %s",
			gecoFile: mockGeco{
				addTemplate: func(gotTemplate geco.Template) {
					internal.AssertEqualString(t, "name", gotTemplate.Name, wantTemplate.Name)
					internal.AssertEqualString(t, "URL", gotTemplate.URL, wantTemplate.URL)
					internal.AssertEqualString(t, "model", gotTemplate.Model, wantTemplate.Model)

					mapGFiles := map[string]bool{}
					for _, gFile := range gotTemplate.GeneratedFiles {
						mapGFiles[gFile] = true
					}

					for _, wantGFile := range wantTemplate.GeneratedFiles {
						if !mapGFiles[wantGFile] {
							t.Errorf("generated file %s not in .geco file", wantGFile)
						}
					}
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := t.TempDir()
			c := fmt.Sprintf(tt.cmd, path)
			args, _ := shell.Parse(c)
			fg := mockFileGenerator{run: func(got renderer.Options, out io.Writer) ([]string, error) {
				return wantTemplate.GeneratedFiles, nil
			}}
			rootCmd := cmd.NewRootCmd(&fg, &tt.gecoFile)
			rootCmd.SetArgs(args)

			err := rootCmd.Execute()
			if err != nil {
				t.Fatalf("Got %s, expected nil", err.Error())
			}
			if tt.gecoFile.NCall != 1 {
				t.Error("expected to be called once")
			}
		})
	}
}

func TestInitErrors(t *testing.T) {
	tests := []errorTestCase{
		{
			name:       "not enough arguments",
			cmd:        "init",
			wantErrMsg: "accepts 2 arg(s), received 0",
		},
	}

	runErrorTestCases(t, tests)
}
