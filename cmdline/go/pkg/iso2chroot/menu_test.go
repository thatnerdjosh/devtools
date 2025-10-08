package iso2chroot

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"thatnerdjosh.com/devtools/pkg/tui"
)

func TestRegisterMenuSelection(t *testing.T) {
	dir := t.TempDir()

	files := []string{"b.iso", "a.iso"}
	for _, name := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(""), 0o644); err != nil {
			t.Fatalf("write file %s: %v", name, err)
		}
	}

	manager := NewManager(dir)

	input := strings.NewReader("1\n2\nq\n")
	var output bytes.Buffer

	menu := tui.NewMenu("iso2chroot — TUI", input, &output)
	RegisterMenu(menu, manager)

	if err := menu.Run(); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	got := output.String()
	if !strings.Contains(got, "Selected ISO: b.iso") {
		t.Fatalf("output does not contain selected ISO: %q", got)
	}
}
