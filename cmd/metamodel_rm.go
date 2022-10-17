/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// metamodelRmCmd represents the rm command
var metamodelRmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove a metamodel",
	Long:  `Remove a metamodel.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("rm called")
	},
}

func init() {
	metamodelCmd.AddCommand(metamodelRmCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// metamodelRmCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// metamodelRmCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
