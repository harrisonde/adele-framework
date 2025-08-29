package helpers

import "os"

// Get environment variable or return default if the value is an empty string.
// Example:
//
//	Helpers.Getenv("ADELE_API_ADDR", "localhost")
func (h *Helpers) Getenv(key string, defaultValue ...string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

// Ensure that a file at the given path exists. If it doesn't, it attempts to create
// the file.
// Example:
//
//	Helpers.CreateFileIfNotExist(path)
func (h *Helpers) CreateFileIfNotExist(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()
	}
	return nil
}

// Ensure that a specific directory exists at the given path. If the directory
// is absent, it proceeds to create it with predefined permissions. This function
// is useful in scenarios where you need to guarantee that a directory is present
// before performing operations that require its existence. A directory that is
// created will have octal value allows the owner to read, write, and execute files
// within the directory, while the group and others can only read and execute, not
// alter the content.
// Example:
//
//	Helpers.CreateDirIfNotExist(path)
func (h *Helpers) CreateDirIfNotExist(path string) error {
	const mode = 0755
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, mode)
		if err != nil {
			return err
		}
	}
	return nil
}
