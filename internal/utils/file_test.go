package utils

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
)

func TestDirExists(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name     string
		setup    func() string
		expected bool
	}{
		{
			name: "Existing directory",
			setup: func() string {
				return tempDir
			},
			expected: true,
		},
		{
			name: "Non-existing directory",
			setup: func() string {
				return filepath.Join(tempDir, "nonexistent")
			},
			expected: false,
		},
		{
			name: "File instead of directory",
			setup: func() string {
				filePath := filepath.Join(tempDir, "testfile.txt")
				os.WriteFile(filePath, []byte("test"), 0644)
				return filePath
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup()
			result := DirExists(path)
			if result != tt.expected {
				t.Errorf("DirExists(%q) = %v, want %v", path, result, tt.expected)
			}
		})
	}
}

func TestCopyRecursiveSync(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(srcDir string)
		wantErr bool
	}{
		{
			name: "Copy single file",
			setup: func(srcDir string) {
				os.WriteFile(filepath.Join(srcDir, "file.txt"), []byte("content"), 0644)
			},
			wantErr: false,
		},
		{
			name: "Copy directory with files",
			setup: func(srcDir string) {
				os.WriteFile(filepath.Join(srcDir, "file1.txt"), []byte("content1"), 0644)
				os.WriteFile(filepath.Join(srcDir, "file2.txt"), []byte("content2"), 0644)
			},
			wantErr: false,
		},
		{
			name: "Copy nested directories",
			setup: func(srcDir string) {
				os.MkdirAll(filepath.Join(srcDir, "subdir1", "subdir2"), 0755)
				os.WriteFile(filepath.Join(srcDir, "subdir1", "file1.txt"), []byte("content1"), 0644)
				os.WriteFile(filepath.Join(srcDir, "subdir1", "subdir2", "file2.txt"), []byte("content2"), 0644)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srcDir := filepath.Join(t.TempDir(), "src")
			destDir := filepath.Join(t.TempDir(), "dest")
			os.MkdirAll(srcDir, 0755)

			tt.setup(srcDir)

			err := CopyRecursiveSync(srcDir, destDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("CopyRecursiveSync() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUnzipFile(t *testing.T) {
	tempDir := t.TempDir()
	zipPath := filepath.Join(tempDir, "test.zip")
	zipFile, err := os.Create(zipPath)
	if err != nil {
		t.Fatal(err)
	}

	zipWriter := zip.NewWriter(zipFile)
	fileWriter, _ := zipWriter.Create("test.txt")
	fileWriter.Write([]byte("test content"))
	zipWriter.Close()
	zipFile.Close()

	destDir := filepath.Join(tempDir, "extracted")
	err = UnzipFile(zipPath, destDir)
	if err != nil {
		t.Fatalf("UnzipFile() error = %v", err)
	}

	extractedFile := filepath.Join(destDir, "test.txt")
	content, err := os.ReadFile(extractedFile)
	if err != nil {
		t.Fatalf("Failed to read extracted file: %v", err)
	}

	if string(content) != "test content" {
		t.Errorf("Content mismatch: got %q, want %q", string(content), "test content")
	}
}
