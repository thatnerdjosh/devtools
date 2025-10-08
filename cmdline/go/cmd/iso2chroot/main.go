package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"thatnerdjosh.com/devtools/pkg/iso2chroot"
	"thatnerdjosh.com/devtools/pkg/tui"
)

const (
	defaultISOdir = "/var/lib/libvirt/isos"
	defaultSrcDir = "/tmp/iso2chroot"
)

func main() {
	flagSet := flag.NewFlagSet("iso2chroot", flag.ContinueOnError)
	flagSet.SetOutput(os.Stderr)

	dir := flagSet.String("dir", defaultISOdir, "Directory containing ISO images")
	src := flagSet.String("src", defaultSrcDir, "Directory to mount the selected ISO into")
	experimentalTUI := flagSet.Bool("experimental-tui", false, "Launch the experimental TUI interface")

	flagSet.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: %s [flags] <command> [args]

Commands:
    list            List available ISO images (default)
    select <index>  Print the ISO identified by its numeric index
    create <index>  Mount the ISO for chroot preparation

Flags:
`, flagSet.Name())
		flagSet.PrintDefaults()
		fmt.Fprint(os.Stderr, `
Examples:
    iso2chroot list
    iso2chroot select 2
    iso2chroot create 1
    iso2chroot --dir /path/to/isos --src /tmp/build-root create 2
`)
	}

	if err := flagSet.Parse(os.Args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			os.Exit(0)
		}
		os.Exit(2)
	}

	manager := iso2chroot.NewManager(*dir)

	if *experimentalTUI {
		if len(flagSet.Args()) > 0 {
			fmt.Fprintln(os.Stderr, "iso2chroot: experimental TUI cannot be combined with CLI commands.")
			os.Exit(2)
		}
		fmt.Fprintln(os.Stderr, "iso2chroot: launching experimental TUI (interface and behavior may change).")
		menu := tui.NewMenu("iso2chroot — TUI (experimental)", os.Stdin, os.Stdout)
		iso2chroot.RegisterMenu(menu, manager)
		if err := menu.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "iso2chroot: %v\n", err)
			os.Exit(1)
		}
		return
	}

	exitCode := iso2chroot.RunCLI(manager, flagSet.Args(), os.Stdout, os.Stderr, iso2chroot.CLIOptions{
		MountDir: *src,
		Stdin:    os.Stdin,
	})
	if exitCode != 0 {
		os.Exit(exitCode)
	}
}
