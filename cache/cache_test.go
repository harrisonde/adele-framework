package cache

import (
	"encoding/json"
	"os"
	"testing"
)

func TestUsesBadger(t *testing.T) {
	// Save original value to restore after test
	original := os.Getenv("CACHE")
	defer os.Setenv("CACHE", original)

	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{
			name:     "CACHE set to badger",
			value:    "badger",
			expected: true,
		},
		{
			name:     "CACHE set to redis",
			value:    "redis",
			expected: false,
		},
		{
			name:     "CACHE empty",
			value:    "",
			expected: false,
		},
		{
			name:     "CACHE set to memory",
			value:    "memory",
			expected: false,
		},
		{
			name:     "CACHE case sensitive",
			value:    "BADGER",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("CACHE", tt.value)
			result := UsesBadger()
			if result != tt.expected {
				t.Errorf("UsesBadger() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUsesRedis(t *testing.T) {
	// Save original values
	originalCache := os.Getenv("CACHE")
	originalSession := os.Getenv("SESSION_TYPE")
	originalQueue := os.Getenv("QUEUE_TYPE")

	defer func() {
		os.Setenv("CACHE", originalCache)
		os.Setenv("SESSION_TYPE", originalSession)
		os.Setenv("QUEUE_TYPE", originalQueue)
	}()

	tests := []struct {
		name        string
		cache       string
		sessionType string
		queueType   string
		expected    bool
	}{
		{
			name:        "CACHE uses redis",
			cache:       "redis",
			sessionType: "",
			queueType:   "",
			expected:    true,
		},
		{
			name:        "SESSION_TYPE uses redis",
			cache:       "",
			sessionType: "redis",
			queueType:   "",
			expected:    true,
		},
		{
			name:        "QUEUE_TYPE uses redis",
			cache:       "",
			sessionType: "",
			queueType:   "redis",
			expected:    true,
		},
		{
			name:        "multiple services use redis",
			cache:       "redis",
			sessionType: "redis",
			queueType:   "",
			expected:    true,
		},
		{
			name:        "no services use redis",
			cache:       "badger",
			sessionType: "cookie",
			queueType:   "memory",
			expected:    false,
		},
		{
			name:        "all variables empty",
			cache:       "",
			sessionType: "",
			queueType:   "",
			expected:    false,
		},
		{
			name:        "case sensitive check",
			cache:       "REDIS",
			sessionType: "",
			queueType:   "",
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("CACHE", tt.cache)
			os.Setenv("SESSION_TYPE", tt.sessionType)
			os.Setenv("QUEUE_TYPE", tt.queueType)

			result := UsesRedis()
			if result != tt.expected {
				t.Errorf("UsesRedis() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestEncode(t *testing.T) {
	tests := []struct {
		name    string
		entry   Entry
		wantErr bool
	}{
		{
			name: "string values",
			entry: Entry{
				"key1": "value1",
				"key2": "value2",
			},
			wantErr: false,
		},
		{
			name: "mixed types",
			entry: Entry{
				"string": "hello",
				"number": float64(42),
				"bool":   true,
			},
			wantErr: false,
		},
		{
			name:    "empty entry",
			entry:   Entry{},
			wantErr: false,
		},
		{
			name: "nested map",
			entry: Entry{
				"nested": map[string]interface{}{
					"inner": "value",
				},
			},
			wantErr: false,
		},
		{
			name: "array values",
			entry: Entry{
				"array": []interface{}{"a", "b", "c"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := Encode(tt.entry)

			if tt.wantErr {
				if err == nil {
					t.Error("Encode() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Encode() unexpected error = %v", err)
				return
			}

			if len(data) == 0 {
				t.Error("Encode() returned empty byte slice")
			}

			// Verify it's valid JSON
			if !isValidJSON(data) {
				t.Error("Encode() did not produce valid JSON")
			}
		})
	}
}

func TestDecode(t *testing.T) {
	// Create valid test data
	validEntry := Entry{
		"test_key": "test_value",
		"number":   float64(123),
	}
	validJSON, _ := Encode(validEntry)

	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name:    "valid JSON",
			data:    validJSON,
			wantErr: false,
		},
		{
			name:    "empty JSON object",
			data:    []byte("{}"),
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			data:    []byte("not json"),
			wantErr: true,
		},
		{
			name:    "empty byte slice",
			data:    []byte{},
			wantErr: true,
		},
		{
			name:    "null JSON",
			data:    []byte("null"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry, err := Decode(tt.data)

			if tt.wantErr {
				if err == nil {
					t.Error("Decode() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Decode() unexpected error = %v", err)
				return
			}

			// Special case: JSON "null" becomes nil Entry, which is valid
			if tt.name == "null JSON" && entry == nil {
				return
			}

			if entry == nil {
				t.Error("Decode() returned nil entry")
			}
		})
	}
}

func TestEncodeDecodeRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		entry Entry
	}{
		{
			name: "simple entry",
			entry: Entry{
				"string_val": "hello",
				"number_val": float64(42),
				"bool_val":   true,
			},
		},
		{
			name: "complex entry",
			entry: Entry{
				"nested": map[string]interface{}{
					"inner": "value",
					"count": float64(5),
				},
				"array": []interface{}{"a", "b", "c"},
			},
		},
		{
			name:  "empty entry",
			entry: Entry{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode
			encoded, err := Encode(tt.entry)
			if err != nil {
				t.Fatalf("Encode() error = %v", err)
			}

			// Decode
			decoded, err := Decode(encoded)
			if err != nil {
				t.Fatalf("Decode() error = %v", err)
			}

			// Compare lengths
			if len(decoded) != len(tt.entry) {
				t.Errorf("Length mismatch: got %d, want %d", len(decoded), len(tt.entry))
			}

			// Compare values (note: JSON numbers become float64)
			for key, originalValue := range tt.entry {
				decodedValue, exists := decoded[key]
				if !exists {
					t.Errorf("Key %s missing from decoded entry", key)
					continue
				}

				// Deep comparison would be more complex, but this covers basic cases
				if !compareValues(originalValue, decodedValue) {
					t.Errorf("Value mismatch for key %s: got %v, want %v", key, decodedValue, originalValue)
				}
			}
		})
	}
}

// Helper functions

func isValidJSON(data []byte) bool {
	var temp interface{}
	return json.Unmarshal(data, &temp) == nil
}

func compareValues(a, b interface{}) bool {
	// Basic comparison - for more complex cases you might need deep comparison
	switch va := a.(type) {
	case map[string]interface{}:
		vb, ok := b.(map[string]interface{})
		if !ok || len(va) != len(vb) {
			return false
		}
		for key, value := range va {
			if !compareValues(value, vb[key]) {
				return false
			}
		}
		return true
	case []interface{}:
		vb, ok := b.([]interface{})
		if !ok || len(va) != len(vb) {
			return false
		}
		for i, value := range va {
			if !compareValues(value, vb[i]) {
				return false
			}
		}
		return true
	default:
		return a == b
	}
}
