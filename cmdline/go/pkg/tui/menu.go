package tui

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

// Action represents a top-level menu action that can optionally install a custom input handler.
type Action func(*Menu) (InputHandler, error)

// InputHandler represents an interaction that consumes subsequent user input.
// Returning nil ends the interaction and returns the menu to its default state.
type InputHandler func(*Menu, string) (InputHandler, error)

// MenuItem represents an entry in the primary menu.
type MenuItem struct {
	Key    string
	Label  string
	Action Action
}

// Menu renders an interactive terminal UI with pluggable menu items.
type Menu struct {
	title      string
	reader     *bufio.Reader
	writer     io.Writer
	content    string
	handler    InputHandler
	items      []MenuItem
	itemLookup map[string]int
	statuses   map[string]string
	statusKeys []string
}

// NewMenu constructs a Menu with the provided reader and writer. A blank title defaults to "Menu".
func NewMenu(title string, reader io.Reader, writer io.Writer) *Menu {
	if title == "" {
		title = "Menu"
	}
	return &Menu{
		title:      title,
		reader:     bufio.NewReader(reader),
		writer:     writer,
		content:    "(content will appear here)",
		itemLookup: make(map[string]int),
		statuses:   make(map[string]string),
	}
}

// AddItem registers a menu item. Keys are matched case-insensitively.
func (m *Menu) AddItem(item MenuItem) {
	key := normalizeKey(item.Key)
	if key == "" {
		return
	}
	if idx, ok := m.itemLookup[key]; ok {
		m.items[idx] = item
		return
	}
	m.itemLookup[key] = len(m.items)
	m.items = append(m.items, item)
}

// SetContent replaces the main content area text.
func (m *Menu) SetContent(content string) {
	m.content = content
}

// AppendContent appends text to the current content, inserting a newline if needed.
func (m *Menu) AppendContent(text string) {
	if m.content == "" {
		m.content = text
		return
	}
	if strings.HasSuffix(m.content, "\n") {
		m.content += text
		return
	}
	m.content += text
}

// SetStatus assigns or updates a status line identified by key. An empty value removes the status.
func (m *Menu) SetStatus(key, value string) {
	key = strings.TrimSpace(key)
	if key == "" {
		return
	}
	if value == "" {
		if _, ok := m.statuses[key]; ok {
			delete(m.statuses, key)
			for i, existing := range m.statusKeys {
				if existing == key {
					m.statusKeys = append(m.statusKeys[:i], m.statusKeys[i+1:]...)
					break
				}
			}
		}
		return
	}
	if _, ok := m.statuses[key]; !ok {
		m.statusKeys = append(m.statusKeys, key)
	}
	m.statuses[key] = value
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

		input := strings.TrimSpace(line)
		lower := strings.ToLower(input)

		if isQuit(lower) {
			m.clearScreen()
			fmt.Fprintln(m.writer, "Goodbye.")
			return nil
		}

		if m.handler != nil {
			next, err := m.handler(m, input)
			if err != nil {
				m.SetContent(fmt.Sprintf("Error: %v", err))
				m.handler = nil
				continue
			}
			m.handler = next
			continue
		}

		if input == "" {
			continue
		}

		if idx, ok := m.itemLookup[lower]; ok {
			item := m.items[idx]
			if item.Action == nil {
				continue
			}
			next, err := item.Action(m)
			if err != nil {
				m.SetContent(fmt.Sprintf("Error: %v", err))
				continue
			}
			m.handler = next
			continue
		}

		m.SetContent(fmt.Sprintf("Unknown choice: %q", input))
	}
}

func (m *Menu) render() {
	m.clearScreen()
	fmt.Fprintln(m.writer, m.title)
	fmt.Fprintln(m.writer, strings.Repeat("-", len(m.title)))
	for _, key := range m.statusKeys {
		if value, ok := m.statuses[key]; ok {
			fmt.Fprintln(m.writer, value)
		}
	}
	if len(m.statusKeys) > 0 {
		fmt.Fprintln(m.writer)
	}

	for _, item := range m.items {
		fmt.Fprintf(m.writer, "[%s] %s\n", item.Key, item.Label)
	}
	fmt.Fprintln(m.writer, "[q] Quit")
	fmt.Fprintln(m.writer)
	fmt.Fprintln(m.writer, "Content:")
	fmt.Fprintln(m.writer, m.content)
	fmt.Fprintln(m.writer)
	fmt.Fprint(m.writer, "> ")
}

func (m *Menu) clearScreen() {
	fmt.Fprint(m.writer, "\033[2J\033[H")
}

func isQuit(input string) bool {
	switch input {
	case "q", "quit", "exit":
		return true
	default:
		return false
	}
}

func normalizeKey(key string) string {
	key = strings.TrimSpace(key)
	if key == "" {
		return ""
	}
	return strings.ToLower(key)
}
