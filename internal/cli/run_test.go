package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunHelp(t *testing.T) {
	for _, args := range [][]string{{"help"}, {"-h"}, {"--help"}, {"migrate", "-h"}} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			var out bytes.Buffer
			orig := commandStdout
			commandStdout = &out
			t.Cleanup(func() { commandStdout = orig })
			if err := Run(args); err != nil {
				t.Fatalf("Run(%q) = %v, want nil", args, err)
			}
			if !strings.HasPrefix(out.String(), "usage: gsxui") {
				t.Fatalf("Run(%q) printed %q, want usage", args, out.String())
			}
		})
	}
}

func TestRunSubcommandFlagHelp(t *testing.T) {
	for _, cmd := range []string{"init", "add", "apply"} {
		if err := Run([]string{cmd, "-h"}); err != nil {
			t.Errorf("Run(%q, -h) = %v, want nil", cmd, err)
		}
	}
}

func TestRunWithoutArgsListsCommands(t *testing.T) {
	err := Run(nil)
	if err == nil {
		t.Fatal("Run(nil) = nil, want usage error")
	}
	for _, cmd := range []string{"init", "list", "add", "apply", "migrate"} {
		if !strings.Contains(err.Error(), "\n  "+cmd+" ") {
			t.Errorf("usage does not describe %s:\n%s", cmd, err)
		}
	}
}
