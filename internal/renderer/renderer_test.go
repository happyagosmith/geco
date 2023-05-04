package renderer_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/happyagosmith/geco/internal"
	"github.com/happyagosmith/geco/internal/renderer"
)

const cT1Template = "this is the parsed template with template1"
const cT1NFTemplate = "this is the parsed template inside the nested folder with template1"
const cT1NFTemplate1 = "this is the parsed template 1 inside the nested folder with template1"

type wantFile = struct {
	path    string
	content string
}
type testCase struct {
	name      string
	opts      renderer.Options
	wantFiles []wantFile
}

type errorTestCase struct {
	name       string
	opts       renderer.Options
	wantErrMsg string
}

func TestApplyAndStdOutput(t *testing.T) {
	t.Run("single file template", func(t *testing.T) {
		opts := renderer.Options{
			TPath: "testdata/template1/template.tpl.md",
			MPath: "testdata/template1/model.yaml",
		}

		want := cT1Template

		b := new(bytes.Buffer)

		r := renderer.Renderer{}
		_, err := r.Run(opts, b)
		if err != nil {
			t.Fatalf("Got %s, expected nil", err.Error())
		}
		got := b.String()
		internal.AssertEqualString(t, "stout", got, want)
	})
}

func TestRendering(t *testing.T) {
	tests := []testCase{
		{
			name: "single file template",
			opts: renderer.Options{
				TPath: "testdata/template1/template.tpl.md",
				MPath: "testdata/template1/model.yaml",
			},
			wantFiles: []wantFile{
				{
					path:    "template.md",
					content: cT1Template,
				},
			},
		},
		{
			name: "folder template",
			opts: renderer.Options{
				TPath: "testdata/template1",
				MPath: "testdata/template1/model.yaml",
			},
			wantFiles: []wantFile{
				{
					path:    "template.md",
					content: cT1Template,
				},
				{
					path:    "nestedFolder/template.md",
					content: cT1NFTemplate,
				},
				{
					path:    "nestedFolder/template1.md",
					content: cT1NFTemplate1,
				},
			},
		},
		{
			name: "merge yaml",
			opts: renderer.Options{
				TPath: "testdata/yaml1",
				MPath: "testdata/yaml1/model.yaml",
			},
			wantFiles: []wantFile{
				{
					path: "merge.yaml",
					content: "key1: v1\n" +
						"key2: v2-overriden\n" +
						"key4:\n" +
						"  key4-1: v4-1-overriden\n" +
						"  key4-2: v4-2\n" +
						"key6:\n" +
						"  - b\n" +
						"  - c\n" +
						"key7:\n" +
						"  - t1: 1\n" +
						"    t2: 2\n" +
						"key3: v3\n" +
						"key5: v5",
				},
			},
		},
		{
			name: "merge not existing yaml",
			opts: renderer.Options{
				TPath: "testdata/yaml2",
				MPath: "testdata/yaml2/model.yaml",
			},
			wantFiles: []wantFile{
				{
					path: "merge.yaml",
					content: "key1: v1\n" +
						"key2: v2\n" +
						"key4:\n" +
						"  key4-1: v4-1\n" +
						"  key4-2: v4-2\n" +
						"key6:\n" +
						"  - a\n" +
						"key7:\n" +
						"  - t1: 1\n" +
						"    t2: 2",
				},
			},
		},
		{
			name: "extract data from comments",
			opts: renderer.Options{
				TPath: "testdata/yaml3",
				MPath: "testdata/yaml3/model.yaml",
			},
			wantFiles: []wantFile{
				{
					path: "README.md",
					content: "### Test1 with title\n" +
						"section description and one scalar parameter\n" +
						"param1.1 desc1.1 default1.1\n" +
						"### Test3 map\n\n" +
						" head comment" +
						"param3.child1 desc3.1 default3.1\n" +
						"param3.child2 desc3.2 default3.2\n" +
						"### Test4 map with items as default values\n\n" +
						"param4 desc4 child1: \"default4.1\"\nchild2: \"default4.2\"\n" +
						"### Test5 empty map\n\n" +
						"param5 desc5 {}\n" +
						"### Test6 array\n\n" +
						"param6 desc6 - key1: value1\n  key2: value2\n" +
						"### Test7 empty array\n\n" +
						"param7 desc7 []\n" +
						"### Test8 annidate map\n\n" +
						"param8.childs desc8  child1: \"default8.1\"\nchild2: \"default8.2\"\n",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.opts.OPath = t.TempDir()

			r := renderer.Renderer{}
			_, err := r.Run(tt.opts, nil)
			if err != nil {
				t.Fatalf("got error %s", err.Error())
			}
			for _, wantFile := range tt.wantFiles {
				got, _ := os.ReadFile(tt.opts.OPath + "/" + wantFile.path)
				internal.AssertEqualString(t, "content of "+wantFile.path, string(got), wantFile.content)
			}
		})
	}
}

func TestErrors(t *testing.T) {

	tests := []errorTestCase{
		{
			name: "template not exists",
			opts: renderer.Options{
				TPath: "testdata/template1/not_exists.tpl.md",
				MPath: "testdata/template1/model.yaml",
			},
			wantErrMsg: "template testdata/template1/not_exists.tpl.md not exists",
		},
		{
			name: "model not exists",
			opts: renderer.Options{
				TPath: "testdata/template1/template.tpl.md",
				MPath: "testdata/template1/not_exists.yaml",
			},
			wantErrMsg: "model testdata/template1/not_exists.yaml not exists",
		},
		{
			name: "output folder not exists",
			opts: renderer.Options{
				TPath: "testdata/template1/template.tpl.md",
				MPath: "testdata/template1/model.yaml",
				OPath: "not_exists/filename",
			},
			wantErrMsg: "output not_exists not exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := renderer.Renderer{}
			_, err := r.Run(tt.opts, nil)
			if err == nil {
				t.Fatalf("Got nil, expected error")
			}
			got := err.Error()

			internal.AssertEqualString(t, "error message", got, tt.wantErrMsg)
		})
	}
}
