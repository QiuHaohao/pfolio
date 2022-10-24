/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/qiuhaohao/pfolio/internal/action"
)

var (
	sortByName       bool
	sortByCreateTime bool
	sortByUpdateTime bool

	descending bool
)

// modelLsCmd represents the list command
var modelLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Ls models",
	Long:  `Ls models.`,
	Run: func(cmd *cobra.Command, args []string) {
		action.NewDefaultAction().ListModels(
			modelSortFlagToSortOrder(
				sortByName, sortByCreateTime, sortByUpdateTime),
			descending)
	},
}

func modelSortFlagToSortOrder(_, byCreateTime, byUpdateTime bool) action.ModelSortOrder {
	if byCreateTime {
		return action.ModelSortOrderByCreateTime
	} else if byUpdateTime {
		return action.ModelSortOrderByUpdateTime
	}
	return action.ModelSortOrderByName
}

func init() {
	modelCmd.AddCommand(modelLsCmd)

	modelLsCmd.Flags().BoolVarP(&sortByName, "sort-by-name", "n", false, "Sort by name(default)")
	modelLsCmd.Flags().BoolVarP(&sortByCreateTime, "sort-by-create-time", "c", false, "Sort by create time")
	modelLsCmd.Flags().BoolVarP(&sortByUpdateTime, "sort-by-update-time", "u", false, "Sort by update time")
	modelLsCmd.MarkFlagsMutuallyExclusive("sort-by-name", "sort-by-create-time", "sort-by-update-time")

	modelLsCmd.Flags().BoolVarP(&descending, "descending", "d", false, "Sort in descending order")

}
