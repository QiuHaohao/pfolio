/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/qiuhaohao/pfolio/internal/cli"
	"github.com/qiuhaohao/pfolio/internal/clock"
	"github.com/qiuhaohao/pfolio/internal/config"
	"github.com/qiuhaohao/pfolio/internal/db"
	"github.com/qiuhaohao/pfolio/internal/editor"
)

// modelUpdateCmd represents the update command
var modelUpdateCmd = &cobra.Command{
	Use:   "update model_name",
	Short: "Update a model",
	Long:  `Update a model.`,
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		m, ok := db.Get().GetModel(name)
		if !ok {
			log.Fatal("Model not found")
		}

		entries, err := editor.EditYamlWithRetry(
			viper.GetString(config.KeyEditor),
			m.Entries,
			func(m db.ModelEntries) error {
				return m.Validate()
			})
		if err != nil {
			log.Fatal(err)
		}

		m.Entries = entries
		m.UpdateTime = clock.Now()

		db.Get().SetModel(name, m)
		db.Persist()

		fmt.Printf("Model %s successfully updated!\n", cli.Highlight(name))
	},
}

func init() {
	modelCmd.AddCommand(modelUpdateCmd)
}
