/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/qiuhaohao/pfolio/internal/action"
)

// modelUpdateCmd represents the update command
var modelUpdateCmd = &cobra.Command{
	Use:   "update model_name",
	Short: "Update a model",
	Long:  `Update a model.`,
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		a := action.NewDefaultAction()
		if err := a.UpdateModel(name); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	modelCmd.AddCommand(modelUpdateCmd)
}
