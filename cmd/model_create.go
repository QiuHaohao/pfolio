/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/QiuHaohao/pfolio/internal/config"
	"github.com/QiuHaohao/pfolio/internal/db"
	"github.com/QiuHaohao/pfolio/internal/editor"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

var name string

// modelCreateCmd represents the create command
var modelCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a model",
	Long:  `Create a model.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if err := db.Get().CheckIsNewModelName(name); err != nil {
			log.Fatal(err)
		}

		var entries db.ModelEntries
		err := config.UnmarshalKey(config.KeyDefaultModel, &entries)
		if err != nil {
			log.Fatal(err)
		}

		entries, err = editor.EditYamlWithRetry(
			viper.GetString(config.KeyEditor),
			entries,
			func(entries db.ModelEntries) error {
				return entries.Validate()
			})
		if err != nil {
			log.Fatal(err)
		}

		if err = db.Get().AddModel(name, entries, false); err != nil {
			log.Fatal(err)
		}

		db.Persist()
	},
}

func init() {
	modelCmd.AddCommand(modelCreateCmd)

	modelCreateCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the new model")
	modelCreateCmd.MarkFlagRequired("name")
}
