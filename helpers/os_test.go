package helpers

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateDirIfNotExist(t *testing.T) {
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "test_dir")

	err := helpers.CreateDirIfNotExist(testPath)
	if err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	// Verify it exists and has correct permissions
	info, err := os.Stat(testPath)
	if err != nil {
		t.Fatal("Directory was not created")
	}
	if !info.IsDir() {
		t.Error("Created path is not a directory")
	}
	if info.Mode().Perm() != 0755 {
		t.Errorf("Wrong permissions: got %o, want 0755", info.Mode().Perm())
	}

	// Test that calling again doesn't error
	err = helpers.CreateDirIfNotExist(testPath)
	if err != nil {
		t.Fatalf("Should not error when directory exists: %v", err)
	}
}

func TestCreateDirIfNotExist_InvalidPath(t *testing.T) {

	// Test with empty path - should error
	err := helpers.CreateDirIfNotExist("")
	if err == nil {
		t.Error("Should fail with empty path")
	}
}

func TestCreateFileIfNotExist(t *testing.T) {
	helpers := &Helpers{}
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test_file.txt")

	// Test creating new file
	err := helpers.CreateFileIfNotExist(testFile)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	// Verify file exists
	info, err := os.Stat(testFile)
	if err != nil {
		t.Fatal("File was not created")
	}
	if info.IsDir() {
		t.Error("Created a directory instead of file")
	}

	// Test that calling again doesn't error
	err = helpers.CreateFileIfNotExist(testFile)
	if err != nil {
		t.Fatalf("Should not error when file exists: %v", err)
	}
}

func TestCreateFileIfNotExist_InvalidPath(t *testing.T) {
	helpers := &Helpers{}

	// Test with invalid path (directory that doesn't exist)
	invalidPath := "/nonexistent/directory/file.txt"
	err := helpers.CreateFileIfNotExist(invalidPath)
	if err == nil {
		t.Error("Should fail when parent directory doesn't exist")
	}
}

func TestCreateFileIfNotExist_EmptyPath(t *testing.T) {
	helpers := &Helpers{}

	// Test with empty path
	err := helpers.CreateFileIfNotExist("")
	if err == nil {
		t.Error("Should fail with empty path")
	}
}

func TestHelpers_Getenv(t *testing.T) {
	helpers := &Helpers{}

	// Test 1: Environment variable exists
	os.Setenv("TEST_VAR", "test_value")
	defer os.Unsetenv("TEST_VAR")

	result := helpers.Getenv("TEST_VAR", "default")
	if result != "test_value" {
		t.Errorf("Expected 'test_value', got '%s'", result)
	}

	// Test 2: Environment variable doesn't exist, use default
	result = helpers.Getenv("NONEXISTENT_VAR", "default_value")
	if result != "default_value" {
		t.Errorf("Expected 'default_value', got '%s'", result)
	}

	// Test 3: Environment variable doesn't exist, no default
	result = helpers.Getenv("NONEXISTENT_VAR")
	if result != "" {
		t.Errorf("Expected empty string, got '%s'", result)
	}

	// Test 4: Environment variable is empty string
	os.Setenv("EMPTY_VAR", "")
	defer os.Unsetenv("EMPTY_VAR")

	result = helpers.Getenv("EMPTY_VAR", "fallback")
	if result != "fallback" {
		t.Errorf("Expected 'fallback' for empty env var, got '%s'", result)
	}

	// Test 5: Multiple default values (should use first)
	result = helpers.Getenv("MISSING_VAR", "first", "second", "third")
	if result != "first" {
		t.Errorf("Expected 'first' from multiple defaults, got '%s'", result)
	}
}

func TestHelpers_Getenv_EmptyKey(t *testing.T) {
	helpers := &Helpers{}

	// Test with empty key
	result := helpers.Getenv("", "default")
	if result != "default" {
		t.Errorf("Expected 'default' for empty key, got '%s'", result)
	}
}
