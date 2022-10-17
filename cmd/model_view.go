/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// modelViewCmd represents the view command
var modelViewCmd = &cobra.Command{
	Use:   "view",
	Short: "View a model",
	Long:  `View a model.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("view called")
	},
}

func init() {
	modelCmd.AddCommand(modelViewCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// modelViewCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// modelViewCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
