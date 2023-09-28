package main

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/harrisonde/adel"
)

var IneritaCommand = &adel.Command{
	Name: "inertia",
	Help: "install classic server-driven web app using inertia",
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
	color.Green("\n  2. install the inertia dependencies:")
	color.Green("    $ npm install")
}
