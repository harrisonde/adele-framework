package adele

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const Version = "v0.0.0"

// Create a new instance of the Adele type using a pointer to Adele with the
// root path of the application as a argument. The new-up is called by project adele's consuming package
// to bootstrap the framework.
func (a *Adele) New(rootPath string) error {

	directories := []string{"data", "handlers", "logs", "jobs", "middleware", "migrations", "public", "resources", "resources/views", "resources/mail", "tmp", "screenshots"}

	err := a.CreateDirectories(rootPath, directories)
	if err != nil {
		return err
	}

	err = a.CreateEnvironmentFile(rootPath)
	if err != nil {
		return err
	}

	err = godotenv.Load(rootPath + "/.env")
	if err != nil {
		return err
	}

	a.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	a.RootPath = rootPath
	a.Version = Version

	return nil
}

// Ensure that a environment file at a specific path exists, creating it if it's missing, and returning
// any errors that may arise.
func (a *Adele) CreateEnvironmentFile(rootPath string) error {
	err := a.CreateFileIfNotExist(fmt.Sprintf("%s", rootPath))
	if err != nil {
		return err
	}
	return nil
}

// Create all nonexistent parent directories
func (a *Adele) CreateDirectories(rootPath string, directories []string) error {
	for _, path := range directories {
		err := a.CreateDirIfNotExist(rootPath + "/" + path)
		if err != nil {
			return err
		}
	}
	return nil
}
