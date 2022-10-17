/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// metamodelCmd represents the metamodel command
var metamodelCmd = &cobra.Command{
	Use:   "metamodel",
	Short: "Manage Metamodels",
	Long:  `Manage Metamodels.`,
}

func init() {
	rootCmd.AddCommand(metamodelCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// metamodelCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// metamodelCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
