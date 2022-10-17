/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// metamodelSyncCmd represents the sync command
var metamodelSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Update the model corresponding to a metamodel",
	Long:  `Update the model corresponding to a metamodel.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("sync called")
	},
}

func init() {
	metamodelCmd.AddCommand(metamodelSyncCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// metamodelSyncCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// metamodelSyncCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
