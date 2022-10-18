/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/qiuhaohao/pfolio/internal/cli"
	"github.com/qiuhaohao/pfolio/internal/db"

	"github.com/spf13/cobra"
)

// modelRmCmd represents the rm command
var modelRmCmd = &cobra.Command{
	Use:   "rm model_name...",
	Short: "Remove a model",
	Long:  `Remove a model.`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, name := range args {
			if !db.Get().ModelNameExists(name) {
				fmt.Printf("Model %s does not exist.\n", cli.Highlight(name))
				continue
			}
			db.Get().RemoveModel(name)
			fmt.Printf("Model %s removed.\n", cli.Highlight(name))
		}

		db.Persist()
	},
}

func init() {
	modelCmd.AddCommand(modelRmCmd)
}
