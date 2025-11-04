package utils

import (
	"reflect"
	"testing"
)

func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		item     string
		expected bool
	}{
		{
			name:     "Item exists in slice",
			slice:    []string{"apple", "banana", "cherry"},
			item:     "banana",
			expected: true,
		},
		{
			name:     "Item does not exist in slice",
			slice:    []string{"apple", "banana", "cherry"},
			item:     "orange",
			expected: false,
		},
		{
			name:     "Empty slice",
			slice:    []string{},
			item:     "apple",
			expected: false,
		},
		{
			name:     "Item is empty string and exists",
			slice:    []string{"apple", "", "cherry"},
			item:     "",
			expected: true,
		},
		{
			name:     "Item is empty string and does not exist",
			slice:    []string{"apple", "banana", "cherry"},
			item:     "",
			expected: false,
		},
		{
			name:     "Case sensitive check",
			slice:    []string{"Apple", "Banana", "Cherry"},
			item:     "apple",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Contains(tt.slice, tt.item)
			if result != tt.expected {
				t.Errorf("Contains(%v, %q) = %v, want %v", tt.slice, tt.item, result, tt.expected)
			}
		})
	}
}

func TestRemoveString(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		item     string
		expected []string
	}{
		{
			name:     "Remove existing item",
			slice:    []string{"apple", "banana", "cherry"},
			item:     "banana",
			expected: []string{"apple", "cherry"},
		},
		{
			name:     "Remove non-existing item",
			slice:    []string{"apple", "banana", "cherry"},
			item:     "orange",
			expected: []string{"apple", "banana", "cherry"},
		},
		{
			name:     "Remove from empty slice",
			slice:    []string{},
			item:     "apple",
			expected: []string{},
		},
		{
			name:     "Remove all occurrences",
			slice:    []string{"apple", "banana", "apple", "cherry", "apple"},
			item:     "apple",
			expected: []string{"banana", "cherry"},
		},
		{
			name:     "Remove empty string",
			slice:    []string{"apple", "", "banana", ""},
			item:     "",
			expected: []string{"apple", "banana"},
		},
		{
			name:     "Single item slice - remove it",
			slice:    []string{"apple"},
			item:     "apple",
			expected: []string{},
		},
		{
			name:     "Single item slice - keep it",
			slice:    []string{"apple"},
			item:     "banana",
			expected: []string{"apple"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RemoveString(tt.slice, tt.item)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("RemoveString(%v, %q) = %v, want %v", tt.slice, tt.item, result, tt.expected)
			}
		})
	}
}

func TestRemoveStringImmutability(t *testing.T) {
	original := []string{"apple", "banana", "cherry"}
	originalCopy := make([]string, len(original))
	copy(originalCopy, original)

	RemoveString(original, "banana")

	if !reflect.DeepEqual(original, originalCopy) {
		t.Errorf("RemoveString modified the original slice: got %v, want %v", original, originalCopy)
	}
}

func BenchmarkContains(b *testing.B) {
	slice := []string{"item1", "item2", "item3", "item4", "item5"}
	item := "item3"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Contains(slice, item)
	}
}

func BenchmarkRemoveString(b *testing.B) {
	slice := []string{"item1", "item2", "item3", "item4", "item5"}
	item := "item3"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RemoveString(slice, item)
	}
}
