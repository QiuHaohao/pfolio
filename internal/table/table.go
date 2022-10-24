package table

import (
	"io"

	"github.com/fatih/color"
	tbl "github.com/rodaine/table"
)

func New(writer io.Writer, columnHeaders ...interface{}) tbl.Table {
	t := tbl.New(columnHeaders...)
	t.WithWriter(writer)
	t.WithHeaderFormatter(color.New(color.Bold, color.FgYellow).Sprintf)
	t.WithFirstColumnFormatter(color.New(color.Bold, color.FgHiWhite).Sprintf)
	return t
}
