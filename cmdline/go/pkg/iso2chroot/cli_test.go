package iso2chroot

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunCLIDefaultList(t *testing.T) {
	dir := t.TempDir()
	files := []string{"b.iso", "a.iso"}
	for _, name := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(""), 0o644); err != nil {
			t.Fatalf("write file %s: %v", name, err)
		}
	}

	manager := NewManager(dir)
	var stdout, stderr bytes.Buffer

	code := RunCLI(manager, nil, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("RunCLI() exit code = %d, want 0", code)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
	got := stdout.String()
	if !strings.Contains(got, " 1. a.iso") || !strings.Contains(got, " 2. b.iso") {
		t.Fatalf("stdout = %q, want entries for a.iso and b.iso", got)
	}
}

func TestRunCLISelect(t *testing.T) {
	dir := t.TempDir()
	files := []string{"b.iso", "a.iso"}
	for _, name := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(""), 0o644); err != nil {
			t.Fatalf("write file %s: %v", name, err)
		}
	}

	manager := NewManager(dir)
	var stdout, stderr bytes.Buffer

	code := RunCLI(manager, []string{"select", "2"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("RunCLI() exit code = %d, want 0", code)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
	if strings.TrimSpace(stdout.String()) != "b.iso" {
		t.Fatalf("stdout = %q, want %q", stdout.String(), "b.iso")
	}
}

func TestRunCLIUnknownCommand(t *testing.T) {
	manager := NewManager(t.TempDir())
	var stdout, stderr bytes.Buffer

	code := RunCLI(manager, []string{"bogus"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("RunCLI() exit code = %d, want 2", code)
	}
	if stderr.Len() == 0 {
		t.Fatal("expected stderr to contain an error message")
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
}

func TestRunCLISelectMissingArgument(t *testing.T) {
	manager := NewManager(t.TempDir())
	var stdout, stderr bytes.Buffer

	code := RunCLI(manager, []string{"select"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("RunCLI() exit code = %d, want 2", code)
	}
	if stderr.Len() == 0 {
		t.Fatal("expected error message on stderr")
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
}
