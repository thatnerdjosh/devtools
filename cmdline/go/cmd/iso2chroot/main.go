package main

import (
	"fmt"
	"os"

	"thatnerdjosh.com/devtools/pkg/iso2chroot"
	"thatnerdjosh.com/devtools/pkg/tui"
)

func main() {
	manager := iso2chroot.NewManager("/var/lib/libvirt/isos")
	menu := tui.NewMenu(manager, os.Stdin, os.Stdout)
	if err := menu.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "iso2chroot: %v\n", err)
		os.Exit(1)
	}
}
