package system

import (
	"fmt"
	"os/exec"
	"regexp"
)

// GetVolumeSerialNumber retrieves the volume serial number
func GetVolumeSerialNumber() (string, error) {
	cmd := exec.Command("cmd", "/c", "vol")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`Serial Number is ([A-F0-9-]+)`)
	match := re.FindStringSubmatch(string(output))
	if len(match) > 1 {
		return match[1], nil
	}

	return "", fmt.Errorf("volume serial number not found")
}

// GetStringSID retrieves the user's SID
func GetStringSID() (string, error) {
	cmd := exec.Command("whoami", "/user")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`\bS-1-\d+-\d+(-\d+)*\b`)
	match := re.FindString(string(output))
	if match != "" {
		// Remove last 5 characters
		if len(match) > 5 {
			return match[:len(match)-5], nil
		}
	}

	return "", fmt.Errorf("SID not found")
}
