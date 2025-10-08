package iso2chroot

import (
	"os"
	"syscall"
	"testing"
)

func TestListISOsOnly(t *testing.T) {
	// NOTE: Unit tests should not necessarily need to touch the filesystem - however, this will do for now.
	// TODO: Calculate this dynamically
	const numISO = 1
	manager := NewManager("fixtures/isos")
	res, err := manager.Load()
	if err != nil {
		t.Error(err)
	}

	if res.Count != numISO {
		t.Errorf("expected %v, got %v", numISO, res.Count)
	}
}

func TestMountISO(t *testing.T) {
	if os.Geteuid() != 0 {
		t.Skip("requires root privileges to mount an ISO")
	}

	manager := NewManager("fixtures/isos")
	if _, err := manager.Load(); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	mountDir := t.TempDir()
	if err := manager.Mount(1, mountDir); err != nil {
		t.Fatalf("Mount failed: %v", err)
	}

	defer func() {
		if err := syscall.Unmount(mountDir, 0); err != nil {
			t.Fatalf("Unmount failed: %v", err)
		}
	}()

	if _, err := os.ReadDir(mountDir); err != nil {
		t.Fatalf("ReadDir failed: %v", err)
	}
}
