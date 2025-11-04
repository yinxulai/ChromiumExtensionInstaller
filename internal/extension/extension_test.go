package extension

import (
"testing"
)

func TestGetExtensionID(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
	}{
		{
			name:     "Simple path",
			filePath: "C:\\test\\path",
		},
		{
			name:     "Empty path",
			filePath: "",
		},
		{
			name:     "Path with backslashes",
			filePath: "C:\\Users\\Test\\AppData\\Roaming\\Extension",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
result := GetExtensionID(tt.filePath)
if len(result) != 32 {
t.Errorf("GetExtensionID() returned ID with length %d, want 32", len(result))
}
for _, c := range result {
if c < 'a' || c > 'p' {
					t.Errorf("GetExtensionID() returned invalid character %c, must be a-p", c)
				}
			}
		})
	}
}

func TestGetExtensionIDConsistency(t *testing.T) {
	path := "C:\\test\\consistency\\path"
	id1 := GetExtensionID(path)
	id2 := GetExtensionID(path)

	if id1 != id2 {
		t.Errorf("GetExtensionID() is not consistent: got %s and %s", id1, id2)
	}
}

func TestGetExtensionIDUniqueness(t *testing.T) {
	paths := []string{
		"C:\\path1",
		"C:\\path2",
		"C:\\different\\path",
		"D:\\another\\location",
	}

	ids := make(map[string]bool)
	for _, path := range paths {
		id := GetExtensionID(path)
		if ids[id] {
			t.Errorf("GetExtensionID() generated duplicate ID %s", id)
		}
		ids[id] = true
	}
}

func TestInstallNonExistentZip(t *testing.T) {
	err := Install("/nonexistent/path/to/extension.zip")
	if err == nil {
		t.Error("Install() should return error for non-existent zip file")
	}
}

func BenchmarkGetExtensionID(b *testing.B) {
	path := "C:\\Users\\Test\\AppData\\Roaming\\BrowserExtensions\\TestExtension"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetExtensionID(path)
	}
}
