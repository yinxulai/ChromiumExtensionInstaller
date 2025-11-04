package browser

import (
	"os"
	"path/filepath"
	"strings"
)

// Browser represents a Chromium-based browser
type Browser struct {
	Name        string
	DisplayName string
	ProfilePath string
	AppPath     string
}

// DetectChromiumBrowsers detects all installed Chromium-based browsers
func DetectChromiumBrowsers() []Browser {
	var browsers []Browser
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return browsers
	}

	localAppData := filepath.Join(homeDir, "AppData", "Local")

	// Define Chromium-based browsers to detect
	browserConfigs := []struct {
		name        string
		displayName string
		profileDir  string
		appDirs     []string
	}{
		{
			name:        "chrome",
			displayName: "Google Chrome",
			profileDir:  filepath.Join(localAppData, "Google", "Chrome", "User Data"),
			appDirs: []string{
				"C:\\Program Files\\Google\\Chrome\\Application",
				"C:\\Program Files (x86)\\Google\\Chrome\\Application",
				filepath.Join(localAppData, "Google", "Chrome", "Application"),
			},
		},
		{
			name:        "edge",
			displayName: "Microsoft Edge",
			profileDir:  filepath.Join(localAppData, "Microsoft", "Edge", "User Data"),
			appDirs: []string{
				"C:\\Program Files\\Microsoft\\Edge\\Application",
				"C:\\Program Files (x86)\\Microsoft\\Edge\\Application",
			},
		},
		{
			name:        "brave",
			displayName: "Brave Browser",
			profileDir:  filepath.Join(localAppData, "BraveSoftware", "Brave-Browser", "User Data"),
			appDirs: []string{
				"C:\\Program Files\\BraveSoftware\\Brave-Browser\\Application",
				"C:\\Program Files (x86)\\BraveSoftware\\Brave-Browser\\Application",
				filepath.Join(localAppData, "BraveSoftware", "Brave-Browser", "Application"),
			},
		},
		{
			name:        "opera",
			displayName: "Opera",
			profileDir:  filepath.Join(homeDir, "AppData", "Roaming", "Opera Software", "Opera Stable"),
			appDirs: []string{
				"C:\\Program Files\\Opera",
				"C:\\Program Files (x86)\\Opera",
				filepath.Join(localAppData, "Programs", "Opera"),
			},
		},
		{
			name:        "vivaldi",
			displayName: "Vivaldi",
			profileDir:  filepath.Join(localAppData, "Vivaldi", "User Data"),
			appDirs: []string{
				"C:\\Program Files\\Vivaldi\\Application",
				"C:\\Program Files (x86)\\Vivaldi\\Application",
				filepath.Join(localAppData, "Vivaldi", "Application"),
			},
		},
		{
			name:        "chromium",
			displayName: "Chromium",
			profileDir:  filepath.Join(localAppData, "Chromium", "User Data"),
			appDirs: []string{
				"C:\\Program Files\\Chromium\\Application",
				"C:\\Program Files (x86)\\Chromium\\Application",
				filepath.Join(localAppData, "Chromium", "Application"),
			},
		},
	}

	for _, config := range browserConfigs {
		// Check if profile directory exists
		if _, err := os.Stat(config.profileDir); err == nil {
			// Find app directory
			var appPath string
			for _, dir := range config.appDirs {
				if _, err := os.Stat(dir); err == nil {
					appPath = dir
					break
				}
			}

			browsers = append(browsers, Browser{
				Name:        config.name,
				DisplayName: config.displayName,
				ProfilePath: config.profileDir,
				AppPath:     appPath,
			})
		}
	}

	return browsers
}

// GetProfilePaths returns all profile directories for a browser
func GetProfilePaths(browser Browser) ([]string, error) {
	if _, err := os.Stat(browser.ProfilePath); os.IsNotExist(err) {
		return nil, err
	}

	var profiles []string
	defaultProfile := filepath.Join(browser.ProfilePath, "Default")
	if _, err := os.Stat(defaultProfile); err == nil {
		profiles = append(profiles, defaultProfile)
	}

	entries, err := os.ReadDir(browser.ProfilePath)
	if err != nil {
		return profiles, nil
	}

	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), "Profile") {
			profiles = append(profiles, filepath.Join(browser.ProfilePath, entry.Name()))
		}
	}

	return profiles, nil
}
