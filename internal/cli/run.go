package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
)

var commandStdin io.Reader = os.Stdin
var commandStdout io.Writer = os.Stdout

// runCommand is the seam for external processes (go, gsx). Unit tests stub
// it; the real implementation streams output through.
var runCommand = func(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

const usage = `usage: gsxui <command> [args]

commands:
  init     set up gsxui in the current Go module: gsxui.json, the ui package, CSS
  list     list the components and helpers gsxui add accepts
  add      vendor components and their dependencies into the ui package
  apply    switch the project's style and theme to a preset
  migrate  rewrite vendored components (gsxui migrate rtl)

Run 'gsxui <command> -h' for a command's flags.`

// Run dispatches the gsxui subcommands.
func Run(args []string) error {
	if len(args) == 0 {
		return errors.New(usage)
	}
	err := dispatch(args)
	if errors.Is(err, flag.ErrHelp) {
		// The flag set has already printed the command's flags.
		return nil
	}
	return err
}

func dispatch(args []string) error {
	switch args[0] {
	case "help", "-h", "-help", "--help":
		fmt.Fprintln(commandStdout, usage)
		return nil
	case "init":
		return runInit(args[1:])
	case "add":
		return runAdd(args[1:])
	case "apply":
		return runApply(args[1:])
	case "migrate":
		return runMigrate(args[1:])
	case "list":
		return runList(args[1:])
	default:
		return fmt.Errorf("unknown command %q (want init, add, apply, migrate, or list; see gsxui help)", args[0])
	}
}
