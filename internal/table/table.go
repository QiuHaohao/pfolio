package table

import (
	"github.com/fatih/color"
	tbl "github.com/rodaine/table"
)

func New(columnHeaders ...interface{}) tbl.Table {
	t := tbl.New(columnHeaders...)
	t.WithHeaderFormatter(color.New(color.Bold, color.FgYellow).Sprintf)
	t.WithFirstColumnFormatter(color.New(color.Bold, color.FgHiWhite).Sprintf)
	return t
}
