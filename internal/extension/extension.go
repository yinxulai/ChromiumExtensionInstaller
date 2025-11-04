package extension

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/yinxulai/chromium-extension-installer/internal/browser"
	"github.com/yinxulai/chromium-extension-installer/internal/system"
	"github.com/yinxulai/chromium-extension-installer/internal/types"
	"github.com/yinxulai/chromium-extension-installer/internal/utils"
)

// GetExtensionID generates the extension ID from the file path
func GetExtensionID(filePath string) string {
	// Convert string to UTF-16LE
	utf16Bytes := utils.EncodeUTF16LE(filePath)

	// Calculate SHA256 hash
	hash := utils.HashSHA256(utf16Bytes)
	digest := hex.EncodeToString(hash[:])

	// Convert to extension ID format
	extensionID := ""
	for i := 0; i < 32; i++ {
		char := digest[i]
		if char >= '0' && char <= '9' {
			extensionID += string(rune('a' + (char - '0')))
		} else {
			extensionID += string(rune('a' + (char - 'a')))
		}
	}

	return extensionID
}

// Install installs a Chrome extension from a zip file
func Install(zipfilePath string) error {
	tempPath := filepath.Join(os.TempDir(), "tempExtensions")
	if err := os.MkdirAll(tempPath, 0755); err != nil {
		return err
	}
	defer os.RemoveAll(tempPath)

	// Extract zip file
	if err := utils.UnzipFile(zipfilePath, tempPath); err != nil {
		return fmt.Errorf("failed to extract zip: %v", err)
	}

	// Read manifest.json
	manifestPath := filepath.Join(tempPath, "manifest.json")
	manifestData, err := os.ReadFile(manifestPath)
	if err != nil {
		return fmt.Errorf("manifest.json not found: %v", err)
	}

	var manifest types.Manifest
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		return fmt.Errorf("failed to parse manifest.json: %v", err)
	}

	extensionName := manifest.Name

	// Copy extension to AppData
	appDataPath := filepath.Join(os.Getenv("APPDATA"), "BrowserExtensions")
	extensionPath := filepath.Join(appDataPath, extensionName)

	if err := os.MkdirAll(extensionPath, 0755); err != nil {
		return err
	}

	entries, err := os.ReadDir(tempPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		src := filepath.Join(tempPath, entry.Name())
		dest := filepath.Join(extensionPath, entry.Name())
		if err := utils.CopyRecursiveSync(src, dest); err != nil {
			return err
		}
	}

	extensionID := GetExtensionID(extensionPath)

	// Get SID and volume serial number
	sid, err := system.GetStringSID()
	if err != nil {
		return fmt.Errorf("failed to get SID: %v", err)
	}

	volumeSerial, err := system.GetVolumeSerialNumber()
	if err != nil {
		return fmt.Errorf("failed to get volume serial number: %v", err)
	}

	fmt.Printf("Extension ID: %s\n", extensionID)
	fmt.Printf("Volume Serial: %s\n", volumeSerial)
	fmt.Printf("SID: %s\n", sid)

	// Detect all Chromium-based browsers
	browsers := browser.DetectChromiumBrowsers()
	if len(browsers) == 0 {
		return fmt.Errorf("no Chromium-based browsers found")
	}

	fmt.Printf("\nDetected %d Chromium-based browser(s):\n", len(browsers))
	for _, b := range browsers {
		fmt.Printf("  - %s\n", b.DisplayName)
	}
	fmt.Println()

	successCount := 0
	// Install to each detected browser
	for _, b := range browsers {
		fmt.Printf("Installing to %s...\n", b.DisplayName)

		// Get encryption key for this browser
		key, err := browser.GetKey(b)
		if err != nil {
			fmt.Printf("  Warning: failed to get key for %s: %v\n", b.DisplayName, err)
			continue
		}

		// Get profiles for this browser
		profiles, err := browser.GetProfilePaths(b)
		if err != nil {
			fmt.Printf("  Warning: failed to get profiles for %s: %v\n", b.DisplayName, err)
			continue
		}

		if len(profiles) == 0 {
			fmt.Printf("  Warning: no profiles found for %s\n", b.DisplayName)
			continue
		}

		// Update each profile
		profileSuccessCount := 0
		for _, profile := range profiles {
			if err := browser.UpdateProfile(profile, extensionID, extensionPath, key, sid); err != nil {
				fmt.Printf("  Warning: failed to update profile %s: %v\n", profile, err)
			} else {
				profileSuccessCount++
			}
		}

		if profileSuccessCount > 0 {
			fmt.Printf("  ✓ Successfully installed to %d profile(s)\n", profileSuccessCount)
			successCount++
		} else {
			fmt.Printf("  ✗ Failed to install to any profile\n")
		}
	}

	if successCount == 0 {
		return fmt.Errorf("failed to install extension to any browser")
	}

	fmt.Printf("\n✓ Extension installed successfully to %d browser(s).\n", successCount)
	return nil
}

// Uninstall removes a Chrome extension
func Uninstall(extensionName string) error {
	appDataPath := filepath.Join(os.Getenv("APPDATA"), "BrowserExtensions")
	extensionPath := filepath.Join(appDataPath, extensionName)

	if _, err := os.Stat(extensionPath); os.IsNotExist(err) {
		return fmt.Errorf("extension not found")
	}

	extensionID := GetExtensionID(extensionPath)

	// Remove extension files
	if err := os.RemoveAll(extensionPath); err != nil {
		return err
	}

	// Get SID
	sid, err := system.GetStringSID()
	if err != nil {
		return fmt.Errorf("failed to get SID: %v", err)
	}

	// Detect all Chromium-based browsers
	browsers := browser.DetectChromiumBrowsers()
	if len(browsers) == 0 {
		return fmt.Errorf("no Chromium-based browsers found")
	}

	fmt.Printf("Detected %d Chromium-based browser(s):\n", len(browsers))
	for _, b := range browsers {
		fmt.Printf("  - %s\n", b.DisplayName)
	}
	fmt.Println()

	successCount := 0
	// Uninstall from each detected browser
	for _, b := range browsers {
		fmt.Printf("Uninstalling from %s...\n", b.DisplayName)

		// Get encryption key for this browser
		key, err := browser.GetKey(b)
		if err != nil {
			fmt.Printf("  Warning: failed to get key for %s: %v\n", b.DisplayName, err)
			continue
		}

		// Get profiles for this browser
		profiles, err := browser.GetProfilePaths(b)
		if err != nil {
			fmt.Printf("  Warning: failed to get profiles for %s: %v\n", b.DisplayName, err)
			continue
		}

		if len(profiles) == 0 {
			fmt.Printf("  Warning: no profiles found for %s\n", b.DisplayName)
			continue
		}

		// Update each profile
		profileSuccessCount := 0
		for _, profile := range profiles {
			if err := browser.RemoveFromProfile(profile, extensionID, key, sid); err != nil {
				fmt.Printf("  Warning: failed to update profile %s: %v\n", profile, err)
			} else {
				profileSuccessCount++
			}
		}

		if profileSuccessCount > 0 {
			fmt.Printf("  ✓ Successfully uninstalled from %d profile(s)\n", profileSuccessCount)
			successCount++
		} else {
			fmt.Printf("  ✗ Failed to uninstall from any profile\n")
		}
	}

	if successCount == 0 {
		return fmt.Errorf("failed to uninstall extension from any browser")
	}

	fmt.Printf("\n✓ Extension uninstalled successfully from %d browser(s).\n", successCount)
	return nil
}
