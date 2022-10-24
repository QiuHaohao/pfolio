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

// modelCreateCmd represents the create command
var modelCreateCmd = &cobra.Command{
	Use:   "create model_name",
	Short: "Create a model",
	Long:  `Create a model.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		initialView := config.MustGetKey[view.ModelEditView](
			config.KeyDefaultModelEditView)

		a := action.NewDefaultAction()
		if err := a.CreateModel(name, initialView); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	modelCmd.AddCommand(modelCreateCmd)
}
