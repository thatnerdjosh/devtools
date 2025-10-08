package iso2chroot

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

const defaultMountDir = "/tmp/iso2chroot"

var mountFunc = func(isoFile, dstDir string) error {
	cmd := exec.Command("sudo", "mount", "-o", "loop,ro", isoFile, dstDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("mount ISO: %w: %s", err, strings.TrimSpace(string(output)))
	}
	return nil
}

// ISOInfo represents a single ISO entry.
type ISOInfo struct {
	Name string
}

// ListResult holds the formatted presentation of ISO options.
type ListResult struct {
	Display string
	Count   int
}

// Manager encapsulates ISO discovery using slice and map structures.
type Manager struct {
	dir         string
	isoByChoice map[int]ISOInfo
	ordered     []ISOInfo
}

// NewManager constructs a Manager rooted at the provided directory.
func NewManager(dir string) *Manager {
	return &Manager{
		dir:         dir,
		isoByChoice: make(map[int]ISOInfo),
		ordered:     make([]ISOInfo, 0),
	}
}

// Directory returns the configured ISO directory.
func (m *Manager) Directory() string {
	return m.dir
}

// Load refreshes ISO entries, populating internal slices and maps.
func (m *Manager) Load() (ListResult, error) {
	entries, err := os.ReadDir(m.dir)
	if err != nil {
		return ListResult{}, fmt.Errorf("read %s: %w", m.dir, err)
	}

	isoEntries := make([]os.DirEntry, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(strings.ToLower(entry.Name()), ".iso") {
			continue
		}
		isoEntries = append(isoEntries, entry)
	}

	if len(isoEntries) == 0 {
		m.isoByChoice = make(map[int]ISOInfo)
		m.ordered = m.ordered[:0]
		return ListResult{
			Display: fmt.Sprintf("No ISO files found in %s", m.dir),
			Count:   0,
		}, nil
	}

	m.isoByChoice = make(map[int]ISOInfo, len(isoEntries))
	if cap(m.ordered) < len(isoEntries) {
		m.ordered = make([]ISOInfo, 0, len(isoEntries))
	} else {
		m.ordered = m.ordered[:0]
	}

	sort.Slice(isoEntries, func(i, j int) bool {
		return isoEntries[i].Name() < isoEntries[j].Name()
	})

	var b strings.Builder
	for i, entry := range isoEntries {
		info := ISOInfo{Name: entry.Name()}
		m.ordered = append(m.ordered, info)
		index := i + 1
		m.isoByChoice[index] = info
		fmt.Fprintf(&b, "%2d. %s\n", index, info.Name)
	}

	return ListResult{
		Display: b.String(),
		Count:   len(m.ordered),
	}, nil
}

func (m *Manager) Mount(choice int, dstDir string) error {
	iso, err := m.Select(choice)
	if err != nil {
		return err
	}

	if dstDir == "" {
		dstDir = defaultMountDir
	}

	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		return fmt.Errorf("prepare mount dir %s: %w", dstDir, err)
	}

	isoFile := filepath.Join(m.dir, iso.Name)
	return mountFunc(isoFile, dstDir)
}

// Select returns the ISO associated with the provided choice number.
func (m *Manager) Select(choice int) (ISOInfo, error) {
	info, ok := m.isoByChoice[choice]
	if !ok {
		return ISOInfo{}, fmt.Errorf("choice %d not available", choice)
	}
	return info, nil
}

// EntryCount reports the number of cached ISO entries.
func (m *Manager) EntryCount() int {
	return len(m.ordered)
}
