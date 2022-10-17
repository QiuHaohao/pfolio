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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// modelUpdateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// modelUpdateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
