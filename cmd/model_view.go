/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/qiuhaohao/pfolio/internal/action"
)

// modelViewCmd represents the view command
var modelViewCmd = &cobra.Command{
	Use:   "view model_name",
	Short: "View a model",
	Long:  `View a model.`,
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		if err := action.NewDefaultAction().ViewModel(name); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	modelCmd.AddCommand(modelViewCmd)
}
