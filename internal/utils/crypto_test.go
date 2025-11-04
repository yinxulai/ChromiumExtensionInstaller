package utils

import (
	"encoding/hex"
	"testing"
)

func TestEncodeUTF16LE(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string // hex encoded for easier comparison
	}{
		{
			name:     "ASCII string",
			input:    "hello",
			expected: "680065006c006c006f00",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Chinese characters",
			input:    "你好",
			expected: "604f7d59",
		},
		{
			name:     "Mixed ASCII and Unicode",
			input:    "Hello世界",
			expected: "480065006c006c006f00164e4c75",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EncodeUTF16LE(tt.input)
			resultHex := hex.EncodeToString(result)
			if resultHex != tt.expected {
				t.Errorf("EncodeUTF16LE(%q) = %s, want %s", tt.input, resultHex, tt.expected)
			}
		})
	}
}

func TestGetHMACSHA256(t *testing.T) {
	tests := []struct {
		name     string
		key      []byte
		message  string
		expected string
	}{
		{
			name:     "Simple message",
			key:      []byte("secret"),
			message:  "message",
			expected: "8b5f48702995c1598c573db1e21866a9b825d4a794d169d7060a03605796360b",
		},
		{
			name:     "Empty message",
			key:      []byte("secret"),
			message:  "",
			expected: "f9e66e179b6747ae54108f82f8ade8b3c25d76fd30afde6c395822c530196169",
		},
		{
			name:     "Empty key",
			key:      []byte(""),
			message:  "message",
			expected: "eb08c1f56d5ddee07f7bdf80468083da06b64cf4fac64fe3a90883df5feacae4",
		},
		{
			name:     "Long message",
			key:      []byte("key"),
			message:  "The quick brown fox jumps over the lazy dog",
			expected: "f7bc83f430538424b13298e6aa6fb143ef4d59a14946175997479dbc2d1a3cd8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetHMACSHA256(tt.key, tt.message)
			if result != tt.expected {
				t.Errorf("GetHMACSHA256() = %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestHashSHA256(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected string // hex encoded hash
	}{
		{
			name:     "Simple data",
			data:     []byte("hello"),
			expected: "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
		},
		{
			name:     "Empty data",
			data:     []byte(""),
			expected: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			name:     "Binary data",
			data:     []byte{0x00, 0x01, 0x02, 0x03, 0x04},
			expected: "08bb5e5d6eaac1049ede0893d30ed022b1a4d9b5b48db414871f51c9cb35283d",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HashSHA256(tt.data)
			resultHex := hex.EncodeToString(result[:])
			if resultHex != tt.expected {
				t.Errorf("HashSHA256() = %s, want %s", resultHex, tt.expected)
			}
		})
	}
}

func BenchmarkEncodeUTF16LE(b *testing.B) {
	input := "Hello, 世界! This is a benchmark test."
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		EncodeUTF16LE(input)
	}
}

func BenchmarkGetHMACSHA256(b *testing.B) {
	key := []byte("secret-key")
	message := "This is a test message for HMAC-SHA256"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetHMACSHA256(key, message)
	}
}

func BenchmarkHashSHA256(b *testing.B) {
	data := []byte("This is a test message for SHA256 hashing")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HashSHA256(data)
	}
}
