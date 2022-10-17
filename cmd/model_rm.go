/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// modelRmCmd represents the rm command
var modelRmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove a model",
	Long:  `Remove a model.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("rm called")
	},
}

func init() {
	modelCmd.AddCommand(modelRmCmd)

	modelRmCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the new model")
	modelRmCmd.MarkFlagRequired("name")
}
