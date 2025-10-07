package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

// clearScreen clears the terminal and moves the cursor to the top-left.
func clearScreen() { fmt.Print("\033[2J\033[H") }

// render draws a barebones menu and a content area.
func render(content string) {
    clearScreen()
    fmt.Println("iso2chroot — TUI (barebones)")
    fmt.Println("----------------------------")
    fmt.Println("[1] List ISOs")
    fmt.Println("[q] Quit")
    fmt.Println()
    fmt.Println("Content:")
    fmt.Println(content)
    fmt.Println()
    fmt.Print("> ")
}

// listISOs uses os.ReadDir to show files in the ISO directory.
func listISOs() string {
    const isoDir = "/var/lib/libvirt/isos"
    entries, err := os.ReadDir(isoDir)
    if err != nil {
        return fmt.Sprintf("Error reading %s: %v", isoDir, err)
    }
    if len(entries) == 0 {
        return fmt.Sprintf("No entries found in %s", isoDir)
    }
    var b strings.Builder
    count := 0
    for i, e := range entries {
        // Show all entries; adjust here if you want to filter.
        fmt.Fprintf(&b, "%2d. %s\n", i+1, e.Name())
        count++
    }
    if count == 0 {
        return "No items to display."
    }
    return b.String()
}

func main() {
    // Preserve existing directory read; result is available for future use.
    entries, err := os.ReadDir("/var/lib/libvirt/isos")
    if err != nil {
        // Do not exit; show the error in the TUI content instead.
    }
    _ = entries

    reader := bufio.NewReader(os.Stdin)
    content := "(content will appear here)"
    if err != nil { // from the preserved read
        content = fmt.Sprintf("Error: %v", err)
    }

    for {
        render(content)
        line, _ := reader.ReadString('\n')
        line = strings.TrimSpace(line)
        switch line {
        case "1":
            content = listISOs()
        case "q", "Q", "quit", "exit":
            clearScreen()
            fmt.Println("Goodbye.")
            return
        default:
            if line != "" {
                content = fmt.Sprintf("Unknown choice: %q", line)
            }
        }
    }
}
