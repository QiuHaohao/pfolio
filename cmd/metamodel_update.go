/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// metamodelUpdateCmd represents the update command
var metamodelUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a metamodel",
	Long:  `Update a metamodel.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("update called")
	},
}

func init() {
	metamodelCmd.AddCommand(metamodelUpdateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// metamodelUpdateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// metamodelUpdateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
