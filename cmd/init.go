/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"path/filepath"

	"github.com/happyagosmith/geco/internal/format"
	"github.com/happyagosmith/geco/internal/geco"
	"github.com/happyagosmith/geco/internal/renderer"

	"github.com/spf13/cobra"
)

type Execution struct {
	Name           string
	URL            string
	GeneratedFiles []string
}

type Result struct {
	Templates []Execution
}

func newInitCmd(fg renderer.FilesGenerator, gf GecoFile) *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "init [template] [dest]",
		Long:  `Initialize the template into dest.`,
		Args: func(cmd *cobra.Command, args []string) error {
			return cobra.ExactArgs(2)(cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			o := renderer.Options{}
			o.TPath = args[0]
			o.OPath = args[1]
			
			writeModel(o.TPath, o.OPath)

			o.MPath = o.OPath + "/geco-model.yaml"
			gFiles, err := fg.Run(o, cmd.OutOrStdout())
			if err != nil {
				return err
			}

			

			gf.AddTemplate(geco.Template{
				Name:           filepath.Base(o.TPath),
				URL:            o.TPath,
				GeneratedFiles: gFiles,
				Model:          "./geco-model.yaml",
			})

			return err
		},
	}
}

func writeModel(tPath, oPath string) {
	mtPath := tPath + "/model.yaml"
	mt, err := os.ReadFile(mtPath)
	if err != nil {
		panic(err)
	}

	moPath := oPath + "/geco-model.yaml"
	_, err = os.Stat(moPath)
	if os.IsNotExist(err) {
		err = os.WriteFile(moPath, mt, 0644)
		if err != nil {
			panic(err)
		}
		return
	}

	mo, err := os.ReadFile(moPath)
	if err != nil {
		panic(err)
	}

	ytm, err := format.NewYaml(mt)
	if err != nil {
		panic(err)
	}

	yom, err := format.NewYaml(mo)
	if err != nil {
		panic(err)
	}

	ytm.Merge(yom)
	b, err := ytm.Bytes()
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(moPath, b, 0644)
	if err != nil {
		panic(err)
	}
}
