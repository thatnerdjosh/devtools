package iso2chroot

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

// RunCLI executes the iso2chroot command-line interface against the provided manager.
// It returns a process exit code, allowing callers to exit appropriately.
func RunCLI(manager *Manager, args []string, stdout, stderr io.Writer) int {
	command := "list"
	if len(args) > 0 {
		command = args[0]
		args = args[1:]
	}

	switch command {
	case "list":
		return runList(manager, stdout, stderr)
	case "select":
		return runSelect(manager, args, stdout, stderr)
	case "help", "-h", "--help":
		fmt.Fprintln(stderr, "iso2chroot commands: list (default), select <index>")
		return 0
	default:
		fmt.Fprintf(stderr, "iso2chroot: unknown command %q\n", command)
		fmt.Fprintln(stderr, "Run 'iso2chroot --help' for usage.")
		return 2
	}
}

func runList(manager *Manager, stdout, stderr io.Writer) int {
	result, err := manager.Load()
	if err != nil {
		fmt.Fprintf(stderr, "iso2chroot: %v\n", err)
		return 1
	}
	printWithTrailingNewline(stdout, result.Display)
	return 0
}

func runSelect(manager *Manager, args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "iso2chroot: select requires a numeric index argument.")
		return 2
	}
	if _, err := manager.Load(); err != nil {
		fmt.Fprintf(stderr, "iso2chroot: %v\n", err)
		return 1
	}

	index, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Fprintf(stderr, "iso2chroot: invalid index %q\n", args[0])
		return 2
	}

	iso, err := manager.Select(index)
	if err != nil {
		fmt.Fprintf(stderr, "iso2chroot: %v\n", err)
		return 1
	}

	fmt.Fprintln(stdout, iso.Name)
	return 0
}

func printWithTrailingNewline(w io.Writer, text string) {
	if text == "" {
		fmt.Fprintln(w)
		return
	}
	if strings.HasSuffix(text, "\n") {
		fmt.Fprint(w, text)
		return
	}
	fmt.Fprintln(w, text)
}
