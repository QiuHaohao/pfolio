/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// metamodelViewCmd represents the view command
var metamodelViewCmd = &cobra.Command{
	Use:   "view",
	Short: "View a metamodel",
	Long:  `View a metamodel.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("view called")
	},
}

func init() {
	metamodelCmd.AddCommand(metamodelViewCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// metamodelViewCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// metamodelViewCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
