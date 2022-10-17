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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// modelRmCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// modelRmCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
