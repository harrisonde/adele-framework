package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/harrisonde/adele-framework"
)

var AuthCommand = &adele.Command{
	Name:        "make auth",
	Help:        "install authentication",
	Description: "install full user authentication into your adel application",
	Usage:       "make auth",
	Options:     map[string]string{},
}

func doAuth() error {
	fmt.Printf("Adele authentication")

	checkForDb()

	color.Yellow("\n\nStarting installation")

	dbType := ade.DB.DataType
	tx, err := ade.PopConnect()
	if err != nil {
		exitGracefully(err)
	}
	defer tx.Close()

	upBytes, err := templateFS.ReadFile("templates/migrations/auth_tables." + dbType + ".sql")
	if err != nil {
		exitGracefully(err)
	}

	downBytes := []byte("drop table if exists users cascade; drop table if exists remember_tokens;")
	if err != nil {
		exitGracefully(err)
	}

	color.Green("  Creating migrations...")

	err = ade.CreatePopMigration(upBytes, downBytes, "auth", "sql")
	if err != nil {
		exitGracefully(err)
	}

	color.Green("  Running migrations...")
	err = ade.RunPopMigrations(tx)
	if err != nil {
		exitGracefully(err)
	}

	appDirs := []string{
		"data",
		"handlers",
		"mail",
		"middleware",
		"migrations",
		"views",
	}

	color.Green("  Creating directories...")
	root := ade.RootPath
	for _, path := range appDirs {
		err := ade.CreateDirIfNotExist(root + "/" + path)
		if err != nil {
			return err
		}
	}

	color.Green("  Creating models...")

	err = copyFileFromTemplate("templates/data/user.go.txt", ade.RootPath+"/data/user.go")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFileFromTemplate("templates/data/remember_token.go.txt", ade.RootPath+"/data/remember_token.go")
	if err != nil {
		exitGracefully(err)
	}

	color.Green("  Creating middleware...")

	err = copyFileFromTemplate("templates/middleware/auth.go.txt", ade.RootPath+"/middleware/auth.go")
	if err != nil {
		exitGracefully(err)
	}

	data, err := templateFS.ReadFile("templates/middleware/remember.go.txt")
	if err != nil {
		return err
	}
	fileName := ade.RootPath + "/middleware/remember.go"
	handler := string(data)
	handler = strings.ReplaceAll(handler, "$APPNAME$", ade.AppName)
	err = ioutil.WriteFile(fileName, []byte(handler), 0644)
	if err != nil {
		return err
	}

	color.Green("  Creating handlers...")

	data, err = templateFS.ReadFile("templates/handlers/auth-handlers.go.txt")
	if err != nil {
		return err
	}
	handler = string(data)
	handler = strings.ReplaceAll(handler, "$APPNAME$", ade.AppName)
	fileName = ade.RootPath + "/handlers/auth-handlers.go"
	err = ioutil.WriteFile(fileName, []byte(handler), 0644)
	if err != nil {
		return err
	}

	color.Green("  Creating mail...")

	err = copyFileFromTemplate("templates/mailer/password-reset.html.tmpl", ade.RootPath+"/mail/password-reset.html.tmpl")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFileFromTemplate("templates/mailer/password-reset.plain.tmpl", ade.RootPath+"/mail/password-reset.plain.tmpl")
	if err != nil {
		exitGracefully(err)
	}

	color.Green("  Creating login, forgot, and password reset views..")
	err = copyFileFromTemplate("templates/views/login.jet", ade.RootPath+"/views/login.jet")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFileFromTemplate("templates/views/forgot.jet", ade.RootPath+"/views/forgot.jet")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFileFromTemplate("templates/views/reset-password.jet", ade.RootPath+"/views/reset-password.jet")
	if err != nil {
		exitGracefully(err)
	}

	cmd := exec.Command("go", "get", "github.com/harrisonde/adele-frameworkeee")
	err = cmd.Start()
	if err != nil {
		exitGracefully(err)
	}

	cmd = exec.Command("go", "mod", "tidy")
	err = cmd.Start()
	if err != nil {
		exitGracefully(err)
	}

	color.Yellow("Installation complete, however additional work required:")
	color.Green("  1. add User model to your application:")
	fmt.Printf("    models.go")
	fmt.Printf("\n      type Models struct {")
	fmt.Printf("\n        Users  User")
	fmt.Printf("\n        RememberToken RememberToken")
	fmt.Printf("\n      }")
	fmt.Printf("\n      return Models{")
	fmt.Printf("\n        Users:          User{},")
	fmt.Printf("\n        RememberToken:  RememberToken{},")
	fmt.Printf("\n      }")
	color.Green("\n  2. add middleware and routes to your web routes:")
	fmt.Printf("    routes-web.go")
	fmt.Printf("\n      middleware:")
	fmt.Printf("\n        r.Use(a.App.NoSurf)")
	fmt.Printf("\n        r.Use(a.Middleware.CheckRemember)")
	fmt.Printf("\n      routes:")
	fmt.Printf("\n        r.Get(\"/login\", a.Handlers.Login)")
	fmt.Printf("\n        r.Post(\"/login\", a.Handlers.PostUserLogin)")
	fmt.Printf("\n        r.Get(\"/logout\", a.Handlers.Logout)")
	fmt.Printf("\n")

	return nil
}
