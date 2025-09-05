package helpers

import (
	"regexp"
	"testing"
)

func TestRandomString(t *testing.T) {
	helpers := &Helpers{}

	// Test different lengths
	lengths := []int{0, 1, 8, 16, 32, 100}

	for _, length := range lengths {
		result := helpers.RandomString(length)

		// Check correct length
		if len(result) != length {
			t.Errorf("Expected length %d, got %d", length, len(result))
		}

		// Check contains only valid characters (alphanumeric)
		validPattern := regexp.MustCompile(`^[a-zA-Z0-9]*$`)
		if !validPattern.MatchString(result) {
			t.Errorf("Generated string contains invalid characters: %s", result)
		}
	}
}

func TestRandomString_Uniqueness(t *testing.T) {
	helpers := &Helpers{}

	// Generate multiple strings and check they're different
	strings := make(map[string]bool)
	length := 16
	iterations := 100

	for i := 0; i < iterations; i++ {
		result := helpers.RandomString(length)

		if strings[result] {
			t.Errorf("Generated duplicate string: %s", result)
		}
		strings[result] = true
	}

	// Should have generated unique strings
	if len(strings) != iterations {
		t.Errorf("Expected %d unique strings, got %d", iterations, len(strings))
	}
}

func TestRandomString_EmptyLength(t *testing.T) {
	helpers := &Helpers{}

	result := helpers.RandomString(0)
	if result != "" {
		t.Errorf("Expected empty string for length 0, got: %s", result)
	}
}
