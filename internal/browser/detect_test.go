package browser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectChromiumBrowsers(t *testing.T) {
	browsers := DetectChromiumBrowsers()

	// Should return a slice (may be nil or empty on Linux or systems without browsers)
	// This is acceptable behavior
	
	// Only validate browser fields if we found any browsers
	for _, browser := range browsers {
		if browser.Name == "" {
			t.Error("Browser name is empty")
		}
		if browser.DisplayName == "" {
			t.Error("Browser display name is empty")
		}
		if browser.ProfilePath == "" {
			t.Error("Browser profile path is empty")
		}
	}
}

func TestBrowserStructure(t *testing.T) {
	browser := Browser{
		Name:        "chrome",
		DisplayName: "Google Chrome",
		ProfilePath: "/path/to/profile",
		AppPath:     "/path/to/app",
	}

	if browser.Name != "chrome" {
		t.Errorf("Name = %s, want chrome", browser.Name)
	}
	if browser.DisplayName != "Google Chrome" {
		t.Errorf("DisplayName = %s, want Google Chrome", browser.DisplayName)
	}
}

func TestGetProfilePathsNonExistent(t *testing.T) {
	browser := Browser{
		Name:        "test",
		DisplayName: "Test Browser",
		ProfilePath: "/nonexistent/path",
		AppPath:     "/nonexistent/app",
	}

	_, err := GetProfilePaths(browser)
	if err == nil {
		t.Error("GetProfilePaths() should return error for non-existent path")
	}
}

func TestGetProfilePathsWithTestData(t *testing.T) {
	tempDir := t.TempDir()
	profilePath := filepath.Join(tempDir, "UserData")
	
	defaultProfile := filepath.Join(profilePath, "Default")
	profile1 := filepath.Join(profilePath, "Profile 1")
	
	os.MkdirAll(defaultProfile, 0755)
	os.MkdirAll(profile1, 0755)

	browser := Browser{
		Name:        "test",
		DisplayName: "Test Browser",
		ProfilePath: profilePath,
		AppPath:     tempDir,
	}

	profiles, err := GetProfilePaths(browser)
	if err != nil {
		t.Fatalf("GetProfilePaths() error = %v", err)
	}

	if len(profiles) < 1 {
		t.Errorf("GetProfilePaths() returned %d profiles, want at least 1", len(profiles))
	}
}

func BenchmarkDetectChromiumBrowsers(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DetectChromiumBrowsers()
	}
}
