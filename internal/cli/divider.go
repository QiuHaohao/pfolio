package cli

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

func PrintDivider() {
	w, _, _ := terminal.GetSize(int(os.Stdin.Fd()))
	fmt.Println(strings.Repeat("―", w))
}
