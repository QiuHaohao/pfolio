/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/qiuhaohao/pfolio/internal/cli"
	"github.com/qiuhaohao/pfolio/internal/db"
	"github.com/qiuhaohao/pfolio/internal/table"
)

// modelViewCmd represents the view command
var modelViewCmd = &cobra.Command{
	Use:   "view model_name",
	Short: "View a model",
	Long:  `View a model.`,
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		m, ok := db.Get().GetModel(name)
		if !ok {
			log.Fatal("Model not found")
		}

		cli.PrintDivider()
		fmt.Printf("Model Name: %s", cli.Highlight(name))
		if m.IsDerivedFromMetaModel {
			fmt.Printf("(Derived from metamodel)")
		}
		fmt.Println()
		fmt.Printf("Create Time: %s\n", m.CreateTime.Format(time.RFC822))
		fmt.Printf("Update Time: %s\n", m.UpdateTime.Format(time.RFC822))

		cli.PrintDivider()
		tbl := table.New("Instrument", "Weight", "Percentage", "Equivalents")

		sort.Slice(m.Entries, func(i, j int) bool {
			return m.Entries[i].Weight > m.Entries[j].Weight ||
				(m.Entries[i].Weight == m.Entries[j].Weight &&
					m.Entries[i].InstrumentIdentifier < m.Entries[i].InstrumentIdentifier)
		})

		totalWeight := m.Entries.TotalWeight()

		for _, e := range m.Entries {
			tbl.AddRow(
				e.InstrumentIdentifier,
				fmt.Sprintf("%.2f", e.Weight),
				fmt.Sprintf("%.2f%%", 100*e.Weight/totalWeight),
				strings.Join(e.GetStringsEquivalentInstruments(), ", "))
		}

		tbl.Print()
	},
}

func init() {
	modelCmd.AddCommand(modelViewCmd)
}
