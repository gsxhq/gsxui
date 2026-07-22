package cli

import (
	"fmt"
	"os"
	"os/exec"
)

// runCommand is the seam for external processes (go, gsx). Unit tests stub
// it; the real implementation streams output through.
var runCommand = func(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Run dispatches the gsxui subcommands.
func Run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: gsxui <init|add|list> [args]")
	}
	switch args[0] {
	case "init":
		return runInit(args[1:])
	case "add":
		return runAdd(args[1:])
	case "list":
		return runList(args[1:])
	default:
		return fmt.Errorf("unknown command %q (want init, add, or list)", args[0])
	}
}
