/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/qiuhaohao/pfolio/internal/action"
)

// modelRmCmd represents the rm command
var modelRmCmd = &cobra.Command{
	Use:   "rm model_name...",
	Short: "Remove a model",
	Long:  `Remove a model.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := action.NewDefaultAction().RemoveModels(args...); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	modelCmd.AddCommand(modelRmCmd)
}
