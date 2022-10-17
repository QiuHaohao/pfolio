/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// metamodelListCmd represents the list command
var metamodelListCmd = &cobra.Command{
	Use:   "list",
	Short: "List metamodels",
	Long:  `List metamodels.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list called")
	},
}

func init() {
	metamodelCmd.AddCommand(metamodelListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// metamodelListCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// metamodelListCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
