package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"unicode/utf16"
)

// EncodeUTF16LE encodes a string to UTF-16LE bytes
func EncodeUTF16LE(s string) []byte {
	utf16Codes := utf16.Encode([]rune(s))
	buf := make([]byte, len(utf16Codes)*2)
	for i, code := range utf16Codes {
		binary.LittleEndian.PutUint16(buf[i*2:], code)
	}
	return buf
}

// GetHMACSHA256 calculates HMAC-SHA256
func GetHMACSHA256(key []byte, message string) string {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

// HashSHA256 calculates SHA256 hash
func HashSHA256(data []byte) [32]byte {
	return sha256.Sum256(data)
}
