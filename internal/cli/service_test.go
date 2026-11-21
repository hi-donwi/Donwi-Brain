package cli_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/hi-donwi/Donwi-Brain/internal/cli"
)

func TestServiceCommand_Tree(t *testing.T) {
	root := cli.NewRootCmd()

	// Verify service command is registered
	cmd, _, err := root.Find([]string{"service"})
	if err != nil || cmd.Name() != "service" {
		t.Fatalf("lokal service command not found: %v", err)
	}

	subcommands := []string{"install", "uninstall", "status"}
	for _, sub := range subcommands {
		subCmd, _, err := root.Find([]string{"service", sub})
		if err != nil || subCmd.Name() != sub {
			t.Errorf("lokal service %s subcommand not found: %v", sub, err)
		}
	}
}

func TestServiceStatus_Help(t *testing.T) {
	root := cli.NewRootCmd()
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs([]string{"service", "--help"})

	err := root.Execute()
	if err != nil {
		t.Fatalf("service --help failed: %v", err)
	}

	helpText := out.String()
	if !strings.Contains(helpText, "install") || !strings.Contains(helpText, "uninstall") || !strings.Contains(helpText, "status") {
		t.Errorf("expected subcommands in help output:\n%s", helpText)
	}
}

func TestServiceInstall_ServeModeFlags(t *testing.T) {
	root := cli.NewRootCmd()
	install, _, err := root.Find([]string{"service", "install"})
	if err != nil {
		t.Fatalf("service install not found: %v", err)
	}
	for _, name := range []string{"cloud-only", "local-only"} {
		if install.Flags().Lookup(name) == nil {
			t.Errorf("service install has no --%s flag", name)
		}
	}
}

func TestServiceInstall_RejectsBothServeModes(t *testing.T) {
	root := cli.NewRootCmd()
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs([]string{"service", "install", "--cloud-only", "--local-only"})

	err := root.Execute()
	if err == nil || !strings.Contains(err.Error(), "cannot combine --cloud-only and --local-only") {
		t.Fatalf("expected the combined modes to be refused, got %v", err)
	}
}
