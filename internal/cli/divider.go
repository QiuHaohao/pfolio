package cli

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

func PrintDivider(stdin *os.File) {
	w, _, _ := terminal.GetSize(int(stdin.Fd()))
	fmt.Println(strings.Repeat("―", w))
}
