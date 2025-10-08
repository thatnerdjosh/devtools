package iso2chroot

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

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

	m.isoByChoice = make(map[int]ISOInfo, len(entries))
	if cap(m.ordered) < len(entries) {
		m.ordered = make([]ISOInfo, 0, len(entries))
	} else {
		m.ordered = m.ordered[:0]
	}

	if len(entries) == 0 {
		return ListResult{
			Display: fmt.Sprintf("No entries found in %s", m.dir),
			Count:   0,
		}, nil
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	var b strings.Builder
	for i, entry := range entries {
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
