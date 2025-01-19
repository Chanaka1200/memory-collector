package checker

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// Function to check if the resctrl filesystem is already mounted

func CheckResctrlMount() bool {

	// Check if the /sys/fs/resctrl directory exists

	_, err := os.Stat("/sys/fs/resctrl")
	if err != nil {
		// If the directory doesn't exist, it's not mounted
		return false
	}

	// Alternatively, you could check the mounted filesystems using `mount` command
	output, err := exec.Command("mount").Output()
	if err != nil {
		log.Fatal(err)
	}

	// Check if /sys/fs/resctrl is present in the list of mounted filesystems
	if strings.Contains(string(output), "/sys/fs/resctrl") {
		return true
	}

	return false
}

// Function to mount the resctrl filesystem with optional parameters

func MountResctrl(options string) error {
	// Mount command with options
	cmd := exec.Command("mount", "-t", "resctrl", "resctrl", "-o", options, "/sys/fs/resctrl")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to mount resctrl filesystem: %v", err)
	}
	return nil
}
