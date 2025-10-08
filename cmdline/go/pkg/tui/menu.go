package tui

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"thatnerdjosh.com/devtools/pkg/iso2chroot"
)

// Menu coordinates user interaction for ISO selection.
type Menu struct {
	manager           *iso2chroot.Manager
	reader            *bufio.Reader
	writer            io.Writer
	content           string
	awaitingSelection bool
	selected          iso2chroot.ISOInfo
	hasSelection      bool
}

// NewMenu constructs a Menu with the provided dependencies.
func NewMenu(manager *iso2chroot.Manager, reader io.Reader, writer io.Writer) *Menu {
	return &Menu{
		manager: manager,
		reader:  bufio.NewReader(reader),
		writer:  writer,
		content: "(content will appear here)",
	}
}

// Run loops the menu until a quit command or input stream end.
func (m *Menu) Run() error {
	for {
		m.render()
		line, err := m.reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		if m.handle(strings.TrimSpace(line)) {
			return nil
		}
	}
}

func (m *Menu) handle(input string) bool {
	if m.awaitingSelection {
		m.handleSelection(input)
		return false
	}

	switch strings.ToLower(input) {
	case "":
		// Keep the existing content.
		return false
	case "1":
		result, err := m.manager.Load()
		if err != nil {
			m.awaitingSelection = false
			m.content = fmt.Sprintf("Error: %v", err)
			return false
		}
		if result.Count == 0 {
			m.awaitingSelection = false
			m.content = result.Display
			return false
		}
		m.awaitingSelection = true
		m.content = fmt.Sprintf("%sEnter the number of the ISO to select it, or 'b' to cancel.", result.Display)
		return false
	case "q", "quit", "exit":
		m.clearScreen()
		fmt.Fprintln(m.writer, "Goodbye.")
		return true
	default:
		m.content = fmt.Sprintf("Unknown choice: %q", input)
		return false
	}
}

func (m *Menu) handleSelection(input string) {
	if input == "" {
		m.content = "Enter a number to choose an ISO, or 'b' to cancel."
		return
	}

	lower := strings.ToLower(input)
	if lower == "b" || lower == "back" {
		m.awaitingSelection = false
		m.content = "Selection cancelled."
		return
	}

	choice, err := strconv.Atoi(input)
	if err != nil {
		m.content = fmt.Sprintf("Invalid selection: %q\nEnter a number between 1 and %d, or 'b' to cancel.", input, m.manager.EntryCount())
		return
	}

	iso, err := m.manager.Select(choice)
	if err != nil {
		m.content = fmt.Sprintf("Invalid selection: %q\nEnter a number between 1 and %d, or 'b' to cancel.", input, m.manager.EntryCount())
		return
	}

	m.awaitingSelection = false
	m.selected = iso
	m.hasSelection = true
	m.content = fmt.Sprintf("Selected ISO: %s", iso.Name)
}

func (m *Menu) render() {
	m.clearScreen()
	fmt.Fprintln(m.writer, "iso2chroot — TUI (barebones)")
	fmt.Fprintln(m.writer, "----------------------------")
	fmt.Fprintln(m.writer, "[1] List ISOs")
	fmt.Fprintln(m.writer, "[q] Quit")
	if m.hasSelection {
		fmt.Fprintf(m.writer, "Selected ISO: %s\n", m.selected.Name)
	}
	fmt.Fprintln(m.writer)
	fmt.Fprintln(m.writer, "Content:")
	fmt.Fprintln(m.writer, m.content)
	fmt.Fprintln(m.writer)
	fmt.Fprint(m.writer, "> ")
}

func (m *Menu) clearScreen() {
	fmt.Fprint(m.writer, "\033[2J\033[H")
}
