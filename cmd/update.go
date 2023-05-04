package cmd

import (
	"github.com/happyagosmith/geco/internal/renderer"
	"github.com/spf13/cobra"
)

func newUpdateCmd(fg renderer.FilesGenerator, gf GecoFile) *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "update [dest]",
		Long:  `Update dest.`,
		Args: func(cmd *cobra.Command, args []string) error {
			return cobra.ExactArgs(1)(cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			tpls := gf.GetTemplates()
			for _, tpl := range tpls {
				o := renderer.Options{
					TPath: tpl.URL,
					OPath: args[0],
					MPath: tpl.Model,
				}
				gFiles, err := fg.Run(o, cmd.OutOrStdout())
				if err != nil {
					return err
				}
				tpl.GeneratedFiles = gFiles
				gf.AddTemplate(tpl)
			}
			return nil
		},
	}
}
