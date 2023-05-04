package cmd_test

import (
	"fmt"
	"io"
	"testing"

	"github.com/happyagosmith/geco/cmd"
	"github.com/happyagosmith/geco/internal"
	"github.com/happyagosmith/geco/internal/geco"
	"github.com/happyagosmith/geco/internal/renderer"
	shell "github.com/mattn/go-shellwords"
)

type testCase struct {
	name     string
	cmd      string
	wantOpts renderer.Options
}

type errorTestCase struct {
	name       string
	cmd        string
	wantErrMsg string
}

type mockFileGenerator struct {
	run   func(o renderer.Options, out io.Writer) ([]string, error)
	NCall int
}

func (m *mockFileGenerator) Run(o renderer.Options, out io.Writer) ([]string, error) {
	m.NCall++
	return m.run(o, out)
}

type mockGeco struct {
	addTemplate  func(t geco.Template)
	getTemplates func() []geco.Template
	NCall        int
}

func (m *mockGeco) AddTemplate(t geco.Template) {
	m.addTemplate(t)
	m.NCall++
}

func (m *mockGeco) GetTemplates() []geco.Template {
	m.NCall++

	return m.getTemplates()
}

func runCheckFilesTests(t *testing.T, tests []testCase) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mg := mockGeco{addTemplate: func(gotTemplate geco.Template) { return }}
			fg := mockFileGenerator{run: func(got renderer.Options, out io.Writer) ([]string, error) {
				internal.AssertEqualString(t, "MPath", got.MPath, tt.wantOpts.MPath)
				internal.AssertEqualString(t, "OPath", got.OPath, tt.wantOpts.OPath)
				internal.AssertEqualString(t, "TPath", got.TPath, tt.wantOpts.TPath)
				return nil, nil
			}}
			args, _ := shell.Parse(tt.cmd)
			rootCmd := cmd.NewRootCmd(&fg, &mg)
			rootCmd.SetArgs(args)
			err := rootCmd.Execute()

			if err != nil {
				t.Fatalf("Got %s, expected nil", err.Error())
			}
		})
	}
}

func runErrorTestCases(t *testing.T, tests []errorTestCase) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := fmt.Sprintf(tt.cmd)
			args, _ := shell.Parse(c)
			mg := mockGeco{addTemplate: func(gotTemplate geco.Template) { return }}
			fg := mockFileGenerator{run: func(got renderer.Options, out io.Writer) ([]string, error) {
				return nil, nil
			}}
			rootCmd := cmd.NewRootCmd(&fg, &mg)
			rootCmd.SetArgs(args)

			err := rootCmd.Execute()
			if err == nil {
				t.Fatalf("Got nil, expected error")
			}
			got := err.Error()

			internal.AssertEqualString(t, "error message", got, tt.wantErrMsg)
		})
	}
}
