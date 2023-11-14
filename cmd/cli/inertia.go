package main

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/harrisonde/adele-framework"
)

var IneritaCommand = &adele.Command{
	Name:        "inertia",
	Help:        "install inertia js",
	Description: "install a classic server driven web application using inertia js",
	Usage:       "inertia",
	Options: map[string]string{
		"-s, --skip": "do not install templates",
	},
}

func doInertiaSetup() {
	root := ade.RootPath
	fmt.Printf("Adele Inertia")
	color.Yellow("\n\nStarting installation")
	color.Green("  Creating package.json...")
	data, err := templateFS.ReadFile("templates/resources/inertia/js/package.json")
	if err != nil {
		exitGracefully(err)
	}

	pj := string(data)
	pj = strings.ReplaceAll(pj, "${APP_NAME}", ade.AppName)

	err = copyDataToFile([]byte(pj), root+"/package.json")
	if err != nil {
		exitGracefully(err)
	}

	longOption, _ := GetOption("skip")
	shortOption, _ := GetOption("s")
	if longOption == "skip" || shortOption == "s" {
		color.Green("  option to skip found, directory creation skipped...")
	} else {
		dirs := []string{
			"resources",
			"resources/css",
			"resources/js",
			"resources/js/components",
			"resources/js/components/forms",
			"resources/js/pages",
			"resources/js/pages/account",
		}

		color.Green("  Creating directories...")

		for _, path := range dirs {
			err := ade.CreateDirIfNotExist(root + "/" + path)
			if err != nil {
				color.Yellow(fmt.Sprintf("%s", err))
			}
			color.Green("    " + path)
		}

		files := []string{
			"css/tailwind.css",
			"js/components/forms/LoginForm.vue",
			"js/components/forms/RegistrationForm.vue",
			"js/components/Flash.vue",
			"js/components/Tag.vue",
			"js/pages/account/index.vue",
			"js/pages/account/register.vue",
			"js/pages/404.vue",
			"js/pages/index.vue",
			"js/app.js",
		}

		color.Green("  Copying files...")
		for _, path := range files {

			d, err := templateFS.ReadFile("templates/resources/inertia/" + path)
			if err != nil {
				exitGracefully(err)
			}

			f := string(d)
			err = copyDataToFile([]byte(f), root+"/resources/"+path)
			if err != nil {
				exitGracefully(err)
			}

			color.Green("    " + path)
		}
	}

	color.Green("  Copying handlers...")
	data, err = templateFS.ReadFile("templates/handlers/inertia-handlers.go.txt")
	if err != nil {
		exitGracefully(err)
	}

	ih := string(data)
	ih = strings.ReplaceAll(ih, "${APP_NAME}", ade.AppName)

	err = copyDataToFile([]byte(ih), root+"/handlers/inertia-handlers.go")
	if err != nil {
		exitGracefully(err)
	}

	color.Green("  Copying models...")
	// User
	data, err = templateFS.ReadFile("templates/data/user.go.txt")
	if err != nil {
		exitGracefully(err)
	}

	um := string(data)
	um = strings.ReplaceAll(um, "${APP_NAME}", ade.AppName)

	err = copyDataToFile([]byte(um), root+"/data/user.go")
	if err != nil {
		exitGracefully(err)
	}

	// Token
	data, err = templateFS.ReadFile("templates/data/token.go.txt")
	if err != nil {
		exitGracefully(err)
	}

	um = string(data)
	um = strings.ReplaceAll(um, "${APP_NAME}", ade.AppName)

	err = copyDataToFile([]byte(um), root+"/data/token.go")
	if err != nil {
		exitGracefully(err)
	}

	// Remember Token
	data, err = templateFS.ReadFile("templates/data/remember_token.go.txt")
	if err != nil {
		exitGracefully(err)
	}

	um = string(data)
	um = strings.ReplaceAll(um, "${APP_NAME}", ade.AppName)

	err = copyDataToFile([]byte(um), root+"/data/remember_token.go")
	if err != nil {
		exitGracefully(err)
	}

	color.Yellow("Installation complete, however additional work required:")
	color.Green("  1. add inertia middleware and routes to your web routes:")
	fmt.Printf("    web-routes.go")
	fmt.Printf("\n      middleware:")
	fmt.Printf("\n        r.Use(a.App.InertiaManager.Middleware)")
	fmt.Printf("\n        r.Use(a.App.NoSurf)")
	fmt.Printf("\n      routes:")
	fmt.Printf("\n        r.Get(\"/{page}\", a.Handlers.Inertia)")
	fmt.Printf("\n        r.Get(\"/{page}/{subpage}\", a.Handlers.Inertia)")
	fmt.Printf("\n        r.Post(\"/account/register\", a.Handlers.PostSignUp)")
	color.Green("\n  2. add User models to your application:")
	fmt.Printf("    models.go")
	fmt.Printf("\n      type Models struct {:")
	fmt.Printf("\n        Users  User")
	fmt.Printf("\n      }")
	fmt.Printf("\n      return Models{")
	fmt.Printf("\n        Users:  User{},")
	fmt.Printf("\n      }")
	color.Green("\n  3. install inertia js dependencies:")
	fmt.Printf("    $ npm install ")
	fmt.Printf("\n")
}
