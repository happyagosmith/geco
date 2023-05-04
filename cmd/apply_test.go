package cmd_test

import (
	"io"
	"testing"

	"github.com/happyagosmith/geco/cmd"
	"github.com/happyagosmith/geco/internal"
	"github.com/happyagosmith/geco/internal/geco"
	"github.com/happyagosmith/geco/internal/renderer"
	shell "github.com/mattn/go-shellwords"
)

func TestApplyAndSave(t *testing.T) {
	tests := []testCase{
		{
			name: "single file template",
			cmd:  "apply -t testdata/template1/template.tpl.md -m testdata/template1/model.yaml -o outputFolder",
			wantOpts: renderer.Options{
				TPath: "testdata/template1/template.tpl.md",
				MPath: "testdata/template1/model.yaml",
				OPath: "outputFolder",
			},
		},
		{
			name: "folder template",
			cmd:  "apply -t testdata/template1 -m testdata/template1/model.yaml -o outputFolder",
			wantOpts: renderer.Options{
				TPath: "testdata/template1",
				MPath: "testdata/template1/model.yaml",
				OPath: "outputFolder",
			},
		},
	}

	runCheckFilesTests(t, tests)
}

func TestUpdateAndSave(t *testing.T) {
	type testCase struct {
		name     string
		cmd      string
		mockGeco mockGeco
		wantOpts []renderer.Options
	}

	tests := []testCase{
		{
			name: "folder template",
			cmd:  "apply -m testdata/template1/model.yaml -o outputFolder",
			mockGeco: mockGeco{
				getTemplates: func() []geco.Template {
					return []geco.Template{
						{
							Name:           "template1",
							URL:            "testdata/template1",
							GeneratedFiles: []string{},
						},
						{
							Name:           "template2",
							URL:            "testdata/template2",
							GeneratedFiles: []string{},
						},
					}
				},
			},
			wantOpts: []renderer.Options{
				{
					TPath: "testdata/template1",
					MPath: "testdata/template1/model.yaml",
					OPath: "outputFolder",
				},
				{
					TPath: "testdata/template2",
					MPath: "testdata/template1/model.yaml",
					OPath: "outputFolder",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fg := mockFileGenerator{}
			fg.run = func(got renderer.Options, out io.Writer) ([]string, error) {
				index := fg.NCall -1
				internal.AssertEqualString(t, "MPath", got.MPath, tt.wantOpts[index].MPath)
				internal.AssertEqualString(t, "OPath", got.OPath, tt.wantOpts[index].OPath)
				internal.AssertEqualString(t, "TPath", got.TPath, tt.wantOpts[index].TPath)
				return nil, nil
			}
			args, _ := shell.Parse(tt.cmd)
			rootCmd := cmd.NewRootCmd(&fg, &tt.mockGeco)
			rootCmd.SetArgs(args)
			err := rootCmd.Execute()

			if err != nil {
				t.Fatalf("Got %s, expected nil", err.Error())
			}
		})
	}
}
