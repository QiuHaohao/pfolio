/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// modelCmd represents the model command
var modelCmd = &cobra.Command{
	Use:   "model",
	Short: "Manage models",
	Long:  `Manage models.`,
}

func init() {
	rootCmd.AddCommand(modelCmd)
}
