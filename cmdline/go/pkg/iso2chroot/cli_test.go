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

	code := RunCLI(manager, nil, &stdout, &stderr, CLIOptions{})
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

	code := RunCLI(manager, []string{"select", "2"}, &stdout, &stderr, CLIOptions{})
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

	code := RunCLI(manager, []string{"bogus"}, &stdout, &stderr, CLIOptions{})
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

	code := RunCLI(manager, []string{"select"}, &stdout, &stderr, CLIOptions{})
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

func TestRunCLICreate(t *testing.T) {
	dir := t.TempDir()
	files := []string{"b.iso", "a.iso"}
	for _, name := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(""), 0o644); err != nil {
			t.Fatalf("write file %s: %v", name, err)
		}
	}

	manager := NewManager(dir)
	var stdout, stderr bytes.Buffer

	var (
		originalMount = mountFunc
		mountCalled   bool
		gotISO        string
		gotDir        string
	)
	defer func() { mountFunc = originalMount }()
	targetDir := filepath.Join(dir, "src")
	mountFunc = func(isoFile, dstDir string) error {
		mountCalled = true
		gotISO = isoFile
		gotDir = dstDir
		return nil
	}

	code := RunCLI(manager, []string{"create", "1"}, &stdout, &stderr, CLIOptions{
		MountDir: targetDir,
		Stdin:    bytes.NewBufferString("\n"),
	})
	if code != 0 {
		t.Fatalf("RunCLI() exit code = %d, want 0", code)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
	if !mountCalled {
		t.Fatal("expected mountFunc to be called")
	}
	wantISO := filepath.Join(dir, "a.iso")
	if gotISO != wantISO {
		t.Fatalf("mounted ISO = %q, want %q", gotISO, wantISO)
	}
	if gotDir != targetDir {
		t.Fatalf("mount dir = %q, want %q", gotDir, targetDir)
	}
	output := stdout.String()
	if !strings.Contains(output, "iso2chroot will mount a.iso into "+targetDir) {
		t.Fatalf("stdout = %q, want warning about mounting into %s", output, targetDir)
	}
	if !strings.Contains(output, "Mounted a.iso to "+targetDir) {
		t.Fatalf("stdout = %q, want success message mentioning %s", output, targetDir)
	}
}

func TestRunCLICreateDefaultMountDir(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "only.iso"), []byte(""), 0o644); err != nil {
		t.Fatalf("write iso: %v", err)
	}

	manager := NewManager(dir)
	var stdout, stderr bytes.Buffer

	var (
		originalMount = mountFunc
		gotDir        string
	)
	defer func() { mountFunc = originalMount }()
	mountFunc = func(isoFile, dstDir string) error {
		gotDir = dstDir
		return nil
	}

	code := RunCLI(manager, []string{"create", "1"}, &stdout, &stderr, CLIOptions{
		Stdin: bytes.NewBufferString("\n"),
	})
	if code != 0 {
		t.Fatalf("RunCLI() exit code = %d, want 0", code)
	}
	if gotDir != defaultMountDir {
		t.Fatalf("mount dir = %q, want default %q", gotDir, defaultMountDir)
	}
}

func TestRunCLICreateMissingArgument(t *testing.T) {
	manager := NewManager(t.TempDir())
	var stdout, stderr bytes.Buffer

	code := RunCLI(manager, []string{"create"}, &stdout, &stderr, CLIOptions{})
	if code != 2 {
		t.Fatalf("RunCLI() exit code = %d, want 2", code)
	}
	if stderr.Len() == 0 {
		t.Fatal("expected error message on stderr")
	}
}

func TestRunCLICreateCancelled(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "only.iso"), []byte(""), 0o644); err != nil {
		t.Fatalf("write iso: %v", err)
	}

	manager := NewManager(dir)
	var stdout, stderr bytes.Buffer

	var (
		originalMount = mountFunc
		mountCalled   bool
	)
	defer func() { mountFunc = originalMount }()
	mountFunc = func(isoFile, dstDir string) error {
		mountCalled = true
		return nil
	}

	code := RunCLI(manager, []string{"create", "1"}, &stdout, &stderr, CLIOptions{
		Stdin: bytes.NewBufferString("n\n"),
	})
	if code != 1 {
		t.Fatalf("RunCLI() exit code = %d, want 1", code)
	}
	if mountCalled {
		t.Fatal("expected mount not to be called")
	}
	if !strings.Contains(stderr.String(), "create cancelled") {
		t.Fatalf("stderr = %q, want cancellation notice", stderr.String())
	}
}
