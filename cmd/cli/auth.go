package main

import (
	"github.com/fatih/color"
)

func doAuth() error {

	checkForDb()

	// Create Migrations
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

	downBytes := []byte("drop table if exists users cascade; drop table if exists tokens cascade; drop table if exists remember_tokens;")
	if err != nil {
		exitGracefully(err)
	}

	err = ade.CreatePopMigration(upBytes, downBytes, "auth", "sql")
	if err != nil {
		exitGracefully(err)
	}

	// Run Migrations
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

	// Create Dirs
	root := ade.RootPath
	for _, path := range appDirs {
		err := ade.CreateDirIfNotExist(root + "/" + path)
		if err != nil {
			return err
		}
	}

	err = copyFileFromTemplate("templates/data/user.go.txt", ade.RootPath+"/data/user.go")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFileFromTemplate("templates/data/token.go.txt", ade.RootPath+"/data/token.go")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFileFromTemplate("templates/data/remember_token.go.txt", ade.RootPath+"/data/remember_token.go")
	if err != nil {
		exitGracefully(err)
	}

	// Copy middleware
	err = copyFileFromTemplate("templates/middleware/auth.go.txt", ade.RootPath+"/middleware/auth.go")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFileFromTemplate("templates/middleware/auth-token.go.txt", ade.RootPath+"/middleware/auth-token.go")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFileFromTemplate("templates/middleware/remember.go.txt", ade.RootPath+"/middleware/remember.go")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFileFromTemplate("templates/handlers/auth-handlers.go.txt", ade.RootPath+"/handlers/auth-handlers.go")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFileFromTemplate("templates/handlers/oauth-handlers.go.txt", ade.RootPath+"/handlers/oauth-handlers.go")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFileFromTemplate("templates/mailer/password-reset.html.tmpl", ade.RootPath+"/mail/password-reset.html.tmpl")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFileFromTemplate("templates/mailer/password-reset.plain.tmpl", ade.RootPath+"/mail/password-reset.plain.tmpl")
	if err != nil {
		exitGracefully(err)
	}

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

	color.Yellow(" - users, tokens, and remembers_tokens migrations created and executed.")
	color.Yellow(" - users and token models created.")
	color.Yellow("")
	color.Yellow("Do not forget to add User and Token models in data/models.go, and to add appropriate middleware to your routes!")

	return nil
}
