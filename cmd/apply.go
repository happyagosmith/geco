/*
Copyright © 2023 Happy Smith happyagosmith@gmail.com
*/
package cmd

import (
	"github.com/happyagosmith/geco/internal/renderer"
	"github.com/spf13/cobra"
)

func newApplyCmd(fg renderer.FilesGenerator, gf GecoFile) *cobra.Command {
	o := renderer.Options{}
	applyCmd := &cobra.Command{
		Use:   "apply",
		Short: "apply -t [template] -m [model]",
		Long:  `Apply values to the template`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if o.TPath != "" {
				_, err := fg.Run(o, cmd.OutOrStdout())
				return err
			}
			return nil
		},
	}
	applyCmd.Flags().StringVarP(&o.MPath, "model", "m", "model.yaml", "Path of the model")
	applyCmd.Flags().StringVarP(&o.OPath, "out", "o", "", "Path where to save the rendered template")
	applyCmd.Flags().StringVarP(&o.TPath, "template", "t", "", "Path of the template")

	_ = applyCmd.MarkFlagRequired("model")

	return applyCmd
}
