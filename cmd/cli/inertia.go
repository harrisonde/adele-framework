package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/harrisonde/adel"
)

var IneritaCommand = &adel.Command{
	Name: "inertia",
	Help: "install classic server-driven web app using inertia",
}

func doInertiaSetup() {
	root := ade.RootPath

	err := copyFileFromTemplate("templates/js/package.json", root+"/package.json")
	if err != nil {
		color.Yellow(fmt.Sprintf("%s", err))
	}

	dirs := []string{
		"resources",
		"resources/css",
		"resources/js",
	}

	for _, path := range dirs {
		err := ade.CreateDirIfNotExist(root + "/" + path)
		if err != nil {
			color.Yellow(fmt.Sprintf("%s", err))
		}
	}

	files := []string{
		"resources/css/tailwind.css",
		"resources/js/components/forms/LoginForm.vue",
		"resources/js/components/forms/RegistrationForm.vue",
		"resources/js/components/Flash.vue",
		"resources/js/components/Tag.vue",
		"resources/js/pages/account/index.vue",
		"resources/js/pages/account/register.vue",
		"resources/js/pages/404.vue",
		"resources/js/pages/index.vue",
		"resources/js/app.js",
	}

	for _, path := range files {
		err := copyFileFromTemplate("templates/"+path, root+"/"+path)
		if err != nil {
			color.Yellow(fmt.Sprintf("%s", err))
		}
	}

	color.Yellow("\n- resources js, css, pages, and, components created")
	color.Yellow("\nadd interita middleware and routes to web-routes.go:\n\n")
	color.Yellow("\tr.Use(a.App.InertiaManager.Middleware)\n\tr.Use(a.App.NoSurf)\n\t...")
	color.Yellow("\nupdate web-routes.go routes with:\n\n")
	color.Yellow("\tr.Get(\"/{page}\", a.Handlers.Inertia)\n\tr.Get(\"/{page}/{subpage}\", a.Handlers.Inertia)\n\tr.Post(\"/account/register\", a.Handlers.PostSignUp)\n\t...")
	color.Yellow("\ninstall the inertia dependencies:\n")
	color.Yellow("\t$ npm install\n\n")
}
