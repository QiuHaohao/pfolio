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
	Use:   "view model_name",
	Short: "View a model",
	Long:  `View a model.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("view called")
	},
}

func init() {
	modelCmd.AddCommand(modelViewCmd)
}
