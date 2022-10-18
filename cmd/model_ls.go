/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"sort"
	"time"

	"github.com/spf13/cobra"

	"github.com/qiuhaohao/pfolio/internal/cli"
	"github.com/qiuhaohao/pfolio/internal/db"
	"github.com/qiuhaohao/pfolio/internal/table"
)

var (
	sortByName       bool
	sortByCreateTime bool
	sortByUpdateTime bool

	descending bool
)

type modelWithName struct {
	name  string
	model db.Model
}

// modelLsCmd represents the list command
var modelLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Ls models",
	Long:  `Ls models.`,
	Run: func(cmd *cobra.Command, args []string) {
		cli.PrintDivider()
		tbl := table.New("Model Name", "Create Time", "Update Time")

		modelsWithName := make([]modelWithName, 0)
		for name, m := range db.Get().Models {
			modelsWithName = append(modelsWithName, modelWithName{name: name, model: m})
		}

		var sortFn func(i, j int) bool

		if sortByCreateTime {
			sortFn = func(i, j int) bool {
				return modelsWithName[i].model.CreateTime.Before(modelsWithName[j].model.CreateTime)
			}
		} else if sortByUpdateTime {
			sortFn = func(i, j int) bool {
				return modelsWithName[i].model.UpdateTime.Before(modelsWithName[j].model.UpdateTime)
			}
		} else {
			sortFn = func(i, j int) bool {
				return modelsWithName[i].name < (modelsWithName[j].name)
			}
		}

		if descending {
			ascSortFn := sortFn
			sortFn = func(i, j int) bool {
				return !ascSortFn(i, j)
			}
		}

		sort.Slice(modelsWithName, sortFn)

		for _, m := range modelsWithName {
			tbl.AddRow(m.name, m.model.CreateTime.Format(time.RFC822), m.model.UpdateTime.Format(time.RFC822))
		}

		tbl.Print()
	},
}

func init() {
	modelCmd.AddCommand(modelLsCmd)

	modelLsCmd.Flags().BoolVarP(&sortByName, "sort-by-name", "n", false, "Sort by name(default)")
	modelLsCmd.Flags().BoolVarP(&sortByCreateTime, "sort-by-create-time", "c", false, "Sort by create time")
	modelLsCmd.Flags().BoolVarP(&sortByUpdateTime, "sort-by-update-time", "u", false, "Sort by update time")
	modelLsCmd.MarkFlagsMutuallyExclusive("sort-by-name", "sort-by-create-time", "sort-by-update-time")

	modelLsCmd.Flags().BoolVarP(&descending, "descending", "d", false, "Sort in descending order")

}
