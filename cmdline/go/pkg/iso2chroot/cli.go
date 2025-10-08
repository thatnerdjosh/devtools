package iso2chroot

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// CLIOptions configures RunCLI behavior.
type CLIOptions struct {
	MountDir string
	Stdin    io.Reader
}

// RunCLI executes the iso2chroot command-line interface against the provided manager.
// It returns a process exit code, allowing callers to exit appropriately.
func RunCLI(manager *Manager, args []string, stdout, stderr io.Writer, opts CLIOptions) int {
	mountDir := opts.MountDir
	if mountDir == "" {
		mountDir = defaultMountDir
	}
	stdin := opts.Stdin
	if stdin == nil {
		stdin = os.Stdin
	}

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
	case "create":
		return runCreate(manager, args, stdout, stderr, mountDir, stdin)
	case "help", "-h", "--help":
		fmt.Fprintln(stderr, "iso2chroot commands: list (default), select <index>, create <index>")
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

func runCreate(manager *Manager, args []string, stdout, stderr io.Writer, mountDir string, stdin io.Reader) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "iso2chroot: create requires a numeric index argument.")
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

	targetDir := mountDir
	if targetDir == "" {
		targetDir = defaultMountDir
	}

	reader := bufio.NewReader(stdin)
	fmt.Fprintf(stdout, "iso2chroot will mount %s into %s using sudo.\n", iso.Name, targetDir)
	fmt.Fprintln(stdout, "You may be prompted for your sudo password.")
	fmt.Fprint(stdout, "Press Enter to continue or type 'n' to cancel: ")

	response, readErr := reader.ReadString('\n')
	if readErr != nil && readErr != io.EOF {
		fmt.Fprintf(stderr, "iso2chroot: read confirmation: %v\n", readErr)
		return 1
	}
	choice := strings.TrimSpace(strings.ToLower(response))
	if choice != "" && choice != "y" && choice != "yes" {
		fmt.Fprintln(stderr, "iso2chroot: create cancelled.")
		return 1
	}

	if err := manager.Mount(index, targetDir); err != nil {
		fmt.Fprintf(stderr, "iso2chroot: %v\n", err)
		return 1
	}

	fmt.Fprintf(stdout, "Mounted %s to %s\n", iso.Name, targetDir)
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
