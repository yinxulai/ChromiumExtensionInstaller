package browser

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// GetKey extracts the encryption key from browser's resources.pak
func GetKey(browser Browser) ([]byte, error) {
	if browser.AppPath == "" {
		return nil, fmt.Errorf("browser application path not found for %s", browser.DisplayName)
	}

	entries, err := os.ReadDir(browser.AppPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read browser directory: %v", err)
	}

	var versionDir string
	versionRegex := regexp.MustCompile(`^\d`)
	for _, entry := range entries {
		if entry.IsDir() && versionRegex.MatchString(entry.Name()) {
			versionDir = entry.Name()
			break
		}
	}

	if versionDir == "" {
		return nil, fmt.Errorf("browser version directory not found for %s", browser.DisplayName)
	}

	resourcesPath := filepath.Join(browser.AppPath, versionDir, "resources.pak")
	if _, err := os.Stat(resourcesPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("resources.pak file not found for %s", browser.DisplayName)
	}

	file, err := os.Open(resourcesPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	buffer := make([]byte, fileInfo.Size())
	_, err = file.Read(buffer)
	if err != nil {
		return nil, err
	}

	offset := 0
	version := binary.LittleEndian.Uint32(buffer[offset:])
	offset += 4

	var resourceCount uint32
	if version == 4 {
		resourceCount = binary.LittleEndian.Uint32(buffer[offset:])
		offset += 4
	} else if version == 5 {
		offset += 4 // skip encoding
		resourceCount = uint32(binary.LittleEndian.Uint16(buffer[offset:]))
		offset += 4
	} else {
		return nil, fmt.Errorf("unsupported resources.pak version: %d", version)
	}

	var key []byte
	prevOffset := uint32(0)
	for i := uint32(0); i < resourceCount; i++ {
		currentOffset := binary.LittleEndian.Uint32(buffer[offset+2:])
		offset += 6

		if i > 0 && currentOffset-prevOffset == 64 {
			key = buffer[prevOffset:currentOffset]
			break
		}

		prevOffset = currentOffset
	}

	if key == nil {
		return nil, fmt.Errorf("key not found in resources.pak for %s", browser.DisplayName)
	}

	return key, nil
}
