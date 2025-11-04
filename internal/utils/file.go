package utils

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

// DirExists checks if a directory exists
func DirExists(dirPath string) bool {
	info, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// CopyRecursiveSync copies files and directories recursively
func CopyRecursiveSync(src, dest string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if srcInfo.IsDir() {
		if err := os.MkdirAll(dest, 0755); err != nil {
			return err
		}

		entries, err := os.ReadDir(src)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			srcFile := filepath.Join(src, entry.Name())
			destFile := filepath.Join(dest, entry.Name())
			if err := CopyRecursiveSync(srcFile, destFile); err != nil {
				return err
			}
		}
	} else {
		srcFile, err := os.Open(src)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		destFile, err := os.Create(dest)
		if err != nil {
			return err
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, srcFile)
		return err
	}

	return nil
}

// UnzipFile extracts a zip file to a destination directory
func UnzipFile(zipPath, destPath string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		filePath := filepath.Join(destPath, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		rc, err := file.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}
