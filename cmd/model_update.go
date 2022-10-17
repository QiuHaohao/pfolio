/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// modelUpdateCmd represents the update command
var modelUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a model",
	Long:  `Update a model.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("update called")
	},
}

func init() {
	modelCmd.AddCommand(modelUpdateCmd)

	modelUpdateCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the new model")
	modelUpdateCmd.MarkFlagRequired("name")
}
