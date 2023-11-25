package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"github.com/harrisonde/adele-framework"
	"github.com/mholt/archiver/v3"
)

var NewCommand = &adele.Command{
	Name:        "new",
	Help:        "create a new application",
	Description: "use this command to create a new adele application",
	Usage:       "new <name>",
	Options: map[string]string{
		"-s, -skip": "do not run go mod tidy",
	},
}

var appURL string

func doNew(appName string) {
	appName = strings.ToLower(appName)
	appURL = appName

	// Clean-up the name and sanitize the user controlled input
	if strings.Contains(appName, "/") {
		exploded := strings.SplitAfter(appName, "/")
		appName = exploded[len(exploded)-1]
	}

	// clone the skeleton application
	color.Green("\tCloning repository...")
	_, err := git.PlainClone("./"+appName, false, &git.CloneOptions{
		URL:      "https://github.com/harrisonde/adele.git",
		Progress: os.Stdout,
		Depth:    1,
	})
	if err != nil {
		exitGracefully(err)
	}

	// remove .git dir
	err = os.RemoveAll(fmt.Sprintf("./%s/.git", appName))
	if err != nil {
		exitGracefully(err)
	}

	// create a ready to roll .env
	color.Yellow("\tCreating .env file ...")
	data, err := templateFS.ReadFile("templates/env.txt")

	env := string(data)
	env = strings.ReplaceAll(env, "${APP_NAME}", appName)
	env = strings.ReplaceAll(env, "${KEY}", ade.RandomString(32))

	err = copyDataToFile([]byte(env), fmt.Sprintf("./%s/.env", appName))
	if err != nil {
		exitGracefully(err)
	}

	// create a makefile
	if runtime.GOOS == "windows" {
		source, err := os.Open(fmt.Sprintf("./%s/Makefile.windows", appName))
		if err != nil {
			exitGracefully(err)
		}
		defer source.Close()
		dest, err := os.Create(fmt.Sprintf("./%s/Makefile", appName))
		if err != nil {
			exitGracefully(err)
		}
		defer dest.Close()

		_, err = io.Copy(dest, source)
		if err != nil {
			exitGracefully(err)
		}
	} else {
		source, err := os.Open(fmt.Sprintf("./%s/Makefile.mac", appName))
		if err != nil {
			exitGracefully(err)
		}
		defer source.Close()
		dest, err := os.Create(fmt.Sprintf("./%s/Makefile", appName))
		if err != nil {
			exitGracefully(err)
		}
		defer dest.Close()

		_, err = io.Copy(dest, source)
		if err != nil {
			exitGracefully(err)
		}
	}

	_ = os.Remove("./" + appName + "/Makefile.mac")
	_ = os.Remove("./" + appName + "/Makefile.windows")

	// update the go.mod for the user
	color.Yellow("\tCreating go.mod file ...")
	_ = os.Remove("./" + appName + "go.mod")

	data, err = templateFS.ReadFile("templates/go.mod.txt")
	if err != nil {
		exitGracefully(err)
	}

	mod := string(data)
	mod = strings.ReplaceAll(mod, "${APP_NAME}", appURL)

	err = copyDataToFile([]byte(mod), "./"+appName+"/go.mod")
	if err != nil {
		exitGracefully(err)
	}

	// update the .go files with the proper imports names
	color.Yellow("\tCreating source files ...")
	os.Chdir("./" + appName)
	updateSource()

	// run go mod tidy in the project dir
	color.Yellow("\tRunning go mod tidy ...")

	cmd := exec.Command("go", "get", "github.com/harrisonde/adele-framework")
	err = cmd.Start()
	if err != nil {
		exitGracefully(err)
	}

	cmd = exec.Command("go", "mod", "tidy")
	err = cmd.Start()
	if err != nil {
		exitGracefully(err)
	}

	// what binary do we need?
	color.Yellow("\tStaring to request cli binary ...")
	binary := "adele-framework_"
	if runtime.GOOS == "darwin" {
		binary = binary + "darwin_x86_64.tar.gz"
	} else if runtime.GOOS == "linux" {
		binary = binary + "linux_arm64.tar.gz"
	} else {
		binary = binary + "windows.exe"
	}

	// build up url to repo
	version := ade.Version
	var url string
	if ade.Version == "" {
		url = "https://github.com/harrisonde/adele-framework/releases/latest/download/" + binary
	} else {
		url = "https://github.com/harrisonde/adele-framework/releases/download/" + version + "/" + binary
	}

	// Download
	color.Yellow("\tDownloading cli binary from " + url)
	err = os.Mkdir("./tmp", 0777)
	if err != nil {
		exitGracefully(err)
	}

	path := "./tmp/" + binary   // update this for package use
	out, err := os.Create(path) // update to
	if err != nil {
		exitGracefully(err)
	}

	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		exitGracefully(err)
	}

	if resp.StatusCode != 200 {
		exitGracefully(errors.New("unable to gracefully download the cli go binary, please manually download from https://github.com/harrisonde/adele-framework/releases and unpack in your project's /cmd directory"))
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		exitGracefully(err)
	}

	// extract the archive
	color.Yellow("\tExtracting cli binary and copying to cmd dir ...")
	err = archiver.Extract(path, "cli", "./cmd")
	if err != nil {
		exitGracefully(err)
	}

	// symlink to root of project
	err = os.Symlink("./cmd/cli", "cli")
	if err != nil {
		fmt.Println(err)
	}

	color.Yellow("\tDone creating " + appName)
	color.White("\tgo build something awesome!")

}
