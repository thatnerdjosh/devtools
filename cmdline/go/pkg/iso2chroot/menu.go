package iso2chroot

import (
	"fmt"
	"strconv"
	"strings"

	"thatnerdjosh.com/devtools/pkg/tui"
)

const isoStatusKey = "iso2chroot.selection"

// RegisterMenu wires iso2chroot interactions into the provided TUI menu.
func RegisterMenu(menu *tui.Menu, manager *Manager) {
	menu.SetStatus(isoStatusKey, "Selected ISO: (none)")
	menu.AddItem(tui.MenuItem{
		Key:   "1",
		Label: "List ISOs",
		Action: func(m *tui.Menu) (tui.InputHandler, error) {
			result, err := manager.Load()
			if err != nil {
				m.SetContent(fmt.Sprintf("Error: %v", err))
				return nil, nil
			}
			if result.Count == 0 {
				m.SetContent(result.Display)
				return nil, nil
			}
			m.SetContent(fmt.Sprintf("%sEnter the number of the ISO to select it, or 'b' to cancel.", result.Display))
			return selectionHandler(manager), nil
		},
	})
}

func selectionHandler(manager *Manager) tui.InputHandler {
	return func(menu *tui.Menu, input string) (tui.InputHandler, error) {
		lower := strings.ToLower(strings.TrimSpace(input))
		switch {
		case lower == "":
			menu.SetContent("Enter a number to choose an ISO, or 'b' to cancel.")
			return selectionHandler(manager), nil
		case lower == "b" || lower == "back":
			menu.SetContent("Selection cancelled.")
			return nil, nil
		}

		choice, err := strconv.Atoi(lower)
		if err != nil {
			menu.SetContent(fmt.Sprintf("Invalid selection: %q\nEnter a number between 1 and %d, or 'b' to cancel.", input, manager.EntryCount()))
			return selectionHandler(manager), nil
		}

		iso, err := manager.Select(choice)
		if err != nil {
			menu.SetContent(fmt.Sprintf("Invalid selection: %q\nEnter a number between 1 and %d, or 'b' to cancel.", input, manager.EntryCount()))
			return selectionHandler(manager), nil
		}

		menu.SetStatus(isoStatusKey, fmt.Sprintf("Selected ISO: %s", iso.Name))
		menu.SetContent(fmt.Sprintf("Selected ISO: %s", iso.Name))
		return nil, nil
	}
}
