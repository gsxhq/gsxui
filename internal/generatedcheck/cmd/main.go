package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/gsxhq/gsxui/internal/generatedcheck"
)

func main() {
	root, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "verify-generated:", err)
		os.Exit(1)
	}
	err = generatedcheck.Check(root, func() error {
		command := exec.Command("go", "tool", "gsx", "generate")
		command.Dir = root
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		return command.Run()
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "verify-generated:", err)
		os.Exit(1)
	}
}
