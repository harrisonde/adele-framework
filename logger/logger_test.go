package logger

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestCreateLogger_Defaults(t *testing.T) {
	// Clear relevant env vars
	os.Unsetenv("LOG_FORMAT")
	os.Unsetenv("DEBUG")
	os.Unsetenv("LOG_LEVEL")

	logger := CreateLogger()
	if logger == nil {
		t.Fatal("CreateLogger() returned nil")
	}

	// Should use TextFormatter by default
	if _, ok := logger.Formatter.(*logrus.TextFormatter); !ok {
		t.Error("Default formatter should be TextFormatter")
	}
}

func TestCreateLogger_JSONFormatter(t *testing.T) {
	os.Setenv("LOG_FORMAT", "JSON")
	defer os.Unsetenv("LOG_FORMAT")

	logger := CreateLogger()
	if _, ok := logger.Formatter.(*logrus.JSONFormatter); !ok {
		t.Error("Should use JSONFormatter when LOG_FORMAT=JSON")
	}
}

func TestCreateLogger_TextFormatter(t *testing.T) {
	os.Setenv("LOG_FORMAT", "TEXT")
	defer os.Unsetenv("LOG_FORMAT")

	logger := CreateLogger()
	if _, ok := logger.Formatter.(*logrus.TextFormatter); !ok {
		t.Error("Should use TextFormatter when LOG_FORMAT!=JSON")
	}
}

func TestCreateLogger_DebugLevel(t *testing.T) {
	os.Setenv("DEBUG", "true")
	defer os.Unsetenv("DEBUG")

	logger := CreateLogger()
	if logger.Level != logrus.DebugLevel {
		t.Errorf("Expected DebugLevel, got %v", logger.Level)
	}
}

func TestCreateLogger_DebugFalse(t *testing.T) {
	os.Setenv("DEBUG", "false")
	defer os.Unsetenv("DEBUG")

	logger := CreateLogger()
	if logger.Level == logrus.DebugLevel {
		t.Error("Should not be DebugLevel when DEBUG=false")
	}
}

func TestCreateLogger_DebugOverridesLogLevel(t *testing.T) {
	// DEBUG=true should override LOG_LEVEL
	os.Setenv("DEBUG", "true")
	os.Setenv("LOG_LEVEL", "ERROR")
	defer func() {
		os.Unsetenv("DEBUG")
		os.Unsetenv("LOG_LEVEL")
	}()

	logger := CreateLogger()
	if logger.Level != logrus.DebugLevel {
		t.Errorf("DEBUG=true should override LOG_LEVEL, got %v", logger.Level)
	}
}

func TestGetLogLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected logrus.Level
	}{
		// Valid levels
		{"panic", logrus.PanicLevel},
		{"fatal", logrus.FatalLevel},
		{"error", logrus.ErrorLevel},
		{"warn", logrus.WarnLevel},
		{"warning", logrus.WarnLevel},
		{"info", logrus.InfoLevel},
		{"debug", logrus.DebugLevel},
		{"trace", logrus.TraceLevel},

		// Case insensitive
		{"PANIC", logrus.PanicLevel},
		{"Fatal", logrus.FatalLevel},
		{"ERROR", logrus.ErrorLevel},
		{"WARN", logrus.WarnLevel},
		{"Warning", logrus.WarnLevel},
		{"INFO", logrus.InfoLevel},
		{"Debug", logrus.DebugLevel},
		{"TRACE", logrus.TraceLevel},

		// Default cases
		{"", logrus.InfoLevel},
		{"invalid", logrus.InfoLevel},
		{"unknown", logrus.InfoLevel},
		{"123", logrus.InfoLevel},
	}

	for _, test := range tests {
		result := GetLogLevel(test.input)
		if result != test.expected {
			t.Errorf("GetLogLevel(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestGetLogLevel_EdgeCases(t *testing.T) {
	// Test with whitespace
	result := GetLogLevel("  debug  ")
	// This will likely return InfoLevel since strings.ToLower doesn't trim
	if result != logrus.InfoLevel {
		t.Errorf("GetLogLevel with whitespace should return InfoLevel (default), got %v", result)
	}

	// Test empty string explicitly
	result = GetLogLevel("")
	if result != logrus.InfoLevel {
		t.Errorf("GetLogLevel(\"\") should return InfoLevel, got %v", result)
	}
}
