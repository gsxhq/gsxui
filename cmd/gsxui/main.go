// gsxui installs shadcn-style gsx components into your project by copying
// their source — you own the code. See https://github.com/gsxhq/gsxui.
package main

import (
	"fmt"
	"os"

	"github.com/gsxhq/gsxui/internal/cli"
)

func main() {
	if err := cli.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "gsxui:", err)
		os.Exit(1)
	}
}
