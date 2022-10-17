package cli

import "github.com/fatih/color"

func Highlight(s string) string {
	return color.New(color.FgHiWhite, color.Bold).Sprintf(s)
}
