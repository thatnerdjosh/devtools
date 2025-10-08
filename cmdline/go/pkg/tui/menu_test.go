package tui

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func scriptedInput(commands ...string) io.Reader {
	lines := append([]string{}, commands...)
	lines = append(lines, "")
	return strings.NewReader(strings.Join(lines, "\n"))
}

func TestMenuActionUpdatesContent(t *testing.T) {
	var output bytes.Buffer
	menu := NewMenu("Test Menu", scriptedInput("1", "q"), &output)

	menu.AddItem(MenuItem{
		Key:   "1",
		Label: "Update content",
		Action: func(m *Menu) (InputHandler, error) {
			m.SetContent("content updated")
			return nil, nil
		},
	})

	if err := menu.Run(); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if got := menu.content; got != "content updated" {
		t.Fatalf("content = %q, want %q", got, "content updated")
	}
}

func TestMenuActionWithHandler(t *testing.T) {
	var output bytes.Buffer
	menu := NewMenu("Test Menu", scriptedInput("1", "42", "q"), &output)

	menu.AddItem(MenuItem{
		Key:   "1",
		Label: "Enter handler",
		Action: func(m *Menu) (InputHandler, error) {
			m.SetContent("handler awaiting value")
			return func(m *Menu, input string) (InputHandler, error) {
				if strings.TrimSpace(input) != "42" {
					return nil, nil
				}
				m.SetContent("handler complete")
				return nil, nil
			}, nil
		},
	})

	if err := menu.Run(); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if got := menu.content; got != "handler complete" {
		t.Fatalf("content = %q, want %q", got, "handler complete")
	}

	if menu.handler != nil {
		t.Fatal("handler should be cleared after completion")
	}
}

func TestMenuUnknownChoiceSetsContent(t *testing.T) {
	var output bytes.Buffer
	menu := NewMenu("Test Menu", scriptedInput("z", "q"), &output)

	if err := menu.Run(); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	want := `Unknown choice: "z"`
	if got := menu.content; got != want {
		t.Fatalf("content = %q, want %q", got, want)
	}
}
