/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/qiuhaohao/pfolio/internal/action"
	"github.com/qiuhaohao/pfolio/internal/config"
	"github.com/qiuhaohao/pfolio/internal/view"
)

// metamodelCreateSyfeCmd represents the syfe command
var metamodelCreateSyfeCmd = &cobra.Command{
	Use:   "syfe metamodel_name",
	Short: "Create a Syfe metamodel",
	Long:  `Create a Syfe metamodel.`,
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		initialView := config.MustGetKey[view.SyfeMetamodelEditView](
			config.KeyDefaultInitialSyfeMetamodelEditView)

		a := action.NewDefaultAction()
		if err := a.CreateSyfeMetamodel(name, initialView); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	metamodelCreateCmd.AddCommand(metamodelCreateSyfeCmd)
}
