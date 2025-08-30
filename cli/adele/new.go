package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/cidekar/adele-framework"
	"github.com/cidekar/adele-framework/helpers"

	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
)

var NewCommand = &Command{
	Name:        "new",
	Help:        "create a new application",
	Description: "use this command to create a new adele application",
	Usage:       "new <name>",
	Options: map[string]string{
		"-p=, --path=":    "where to create the new project",
		"-v=, --version=": "specify a version of adele to install",
	},
}

var repositoryRoot = "git.86labs.cloud/harrison"

type CommandNewApplication interface {
	Handle() error
	AddBinary() error
	Clone(string) error
	Sanitize(string) string
	UpdateSource(string) error
	Validate(string) error
	Write(string, string) error
}

var modulePath string

type NewApp struct {
	command CommandNewApplication
}

func NewApplication() *NewApp {
	return &NewApp{}
}

func (c *NewApp) AddBinary() error {

	// run go mod tidy in the project dir
	color.Yellow("\tRunning go mod tidy ...")

	var version string
	if HasOption("version") || HasOption("v") {
		version, _ = GetOption("version")
	} else {
		version = adele.Version
	}

	color.Yellow("\t install adele framework version " + version)
	cmd := exec.Command("go", "get", repositoryRoot+"/adele-framework@"+version)
	err := cmd.Start()
	if err != nil {
		return err
	}

	cmd = exec.Command("go", "mod", "tidy")
	err = cmd.Start()
	if err != nil {
		return err
	}

	var binary string
	// what binary do we need?
	color.Yellow("\tStaring to request cli binary ...")
	if runtime.GOOS == "darwin" {
		binary = "cli_darwin_x86_64"
	} else if runtime.GOOS == "linux" {
		binary = "cli_linux_arm64"
	} else {
		binary = "cli_windows"
	}

	url := "https://" + repositoryRoot + "/adele-framework/releases/download/" + adele.Version + "/" + binary

	// Download
	tmpDirPath := "./tmp"
	color.Yellow("\tDownloading cli binary from " + url)
	if _, err := os.Stat(tmpDirPath); os.IsNotExist(err) {
		err = os.Mkdir(tmpDirPath, 0777)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}

	path := "./cmd/cli" + binary
	out, err := os.Create(path)
	if err != nil {
		return err
	}

	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("unable to gracefully download the cli go binary, please manually download from https://github.com/cidekar/adele-framework/releases and unpack in your project's /cmd directory")
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	err = os.Rename(path, "./cmd/cli")
	if err != nil {
		return err
	}

	err = os.Chmod("./cmd/cli", 0744)
	if err != nil {
		return err
	}

	// symlink to root of project
	err = os.Symlink("./cmd/cli", "cli")
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func (c *NewApp) Clone(appName string) error {

	if flag.Lookup("test.v") != nil {
		color.Yellow("\tRunning go test ... skip Clone")
	} else {
		// clone the skeleton application
		color.Green("\tCloning repository...")

		_, err := git.PlainClone("./"+appName, false, &git.CloneOptions{
			URL:      "https://" + repositoryRoot + "/adele.git",
			Progress: os.Stdout,
			Depth:    1,
		})

		if err != nil {
			return err
		}

		// remove .git dir
		err = os.RemoveAll(fmt.Sprintf("./%s/.git", appName))
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *NewApp) Sanitize(appName string) string {
	name := strings.ToLower(appName)

	if strings.Contains(name, "/") {
		exploded := strings.SplitAfter(name, "/")
		name = exploded[len(exploded)-1]
	}

	return name
}

func (c *NewApp) Validate(arg3 string) error {
	if arg3 == "" {
		return errors.New("you must provide a name for the application")
	}
	return nil
}

func (c *NewApp) UpdateSource(modulePathFromCaller string) error {
	modulePath = modulePathFromCaller
	err := filepath.Walk(".", updateSourceFiles)
	if err != nil {
		return err
	}
	return nil
}

func (c *NewApp) Write(appName, modulePath string) error {

	if HasOption("path") || HasOption("p") {
		p, _ := GetOption("path")
		color.Yellow("\t" + p)
		os.Chdir(p)
	}

	// create a ready-to-roll env file
	color.Yellow("\tCreating .env file ...")
	data, err := templateFS.ReadFile("templates/env.txt")
	if err != nil {
		return err
	}

	env := string(data)
	env = strings.ReplaceAll(env, "${APP_NAME}", appName)

	helpers := helpers.Helpers{}
	env = strings.ReplaceAll(env, "${KEY}", helpers.RandomString(32))

	err = copyDataToFile([]byte(env), fmt.Sprintf("./%s/.env", appName))
	if err != nil {
		return err
	}

	// create a makefile
	if runtime.GOOS == "windows" {
		source, err := os.Open(fmt.Sprintf("./%s/Makefile.windows", appName))
		if err != nil {
			return err
		}
		defer source.Close()
		dest, err := os.Create(fmt.Sprintf("./%s/Makefile", appName))
		if err != nil {
			return err
		}
		defer dest.Close()

		_, err = io.Copy(dest, source)
		if err != nil {
			return err
		}
	} else {
		source, err := os.Open(fmt.Sprintf("./%s/Makefile.mac", appName))
		if err != nil {
			return err
		}
		defer source.Close()
		dest, err := os.Create(fmt.Sprintf("./%s/Makefile", appName))
		if err != nil {
			return err
		}
		defer dest.Close()

		_, err = io.Copy(dest, source)
		if err != nil {
			return err
		}
	}

	_ = os.Remove("./" + appName + "/Makefile.mac")
	_ = os.Remove("./" + appName + "/Makefile.windows")

	// update the go.mod for the user
	color.Yellow("\tCreating go.mod file ...")
	_ = os.Remove("./" + appName + "go.mod")

	data, err = templateFS.ReadFile("templates/go.mod.txt")
	if err != nil {
		return err
	}

	mod := string(data)
	mod = strings.ReplaceAll(mod, "${APP_MODULE_PATH}", modulePath)
	mod = strings.ReplaceAll(mod, "${ADELE_PACKAGE_VERSION}", adele.Version)

	err = copyDataToFile([]byte(mod), "./"+appName+"/go.mod")
	if err != nil {
		return err
	}

	// update the .go files with the proper imports names
	color.Yellow("\tCreating source files ...")
	os.Chdir("./" + appName)

	err = c.UpdateSource(modulePath)
	if err != nil {
		return err
	}

	return nil
}

func (c *NewApp) Handle(appName string) error {
	err := c.Validate(appName)
	if err != nil {
		return (err)
	}

	modulePath := appName
	name := c.Sanitize(appName)

	if HasOption("path") || HasOption("p") {
		p, _ := GetOption("path")
		color.Yellow("\t" + p)
		os.Chdir(p)
	}

	err = c.Clone(name)
	if err != nil {
		return err
	}

	err = c.Write(name, modulePath)
	if err != nil {
		return (err)
	}

	if flag.Lookup("test.v") != nil {
		color.Yellow("\tRunning go test ... skip AddBinary")
		return nil
	} else {
		err = c.AddBinary()
		if err != nil {
			return err
		}
	}

	color.Yellow("\tDone creating " + name)
	color.White("\tgo build something awesome!")

	return nil
}

func updateSourceFiles(path string, fi os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if fi.IsDir() {
		return nil
	}

	matched, err := filepath.Match("*.go", fi.Name())
	if err != nil {
		return err
	}
	if matched {
		read, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		newCont := strings.Replace(string(read), "myapp", modulePath, -1)
		err = os.WriteFile(path, []byte(newCont), 0)
		if err != nil {
			return err
		}
	}

	return nil
}
