package main

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

func doAuth() error {
	// Create Migrations
	dbType := ade.DB.DataType
	fileName := fmt.Sprintf("%d_create_auth_tables", time.Now().UnixMicro())
	upFile := ade.RootPath + "/migrations/" + fileName + ".up.sql"
	downFile := ade.RootPath + "/migrations/" + fileName + ".down.sql"

	err := copyFileFromTemplate("templates/migrations/auth_tables."+dbType+".sql", upFile)
	if err != nil {
		exitGracefully(err)
	}

	err = copyDataToFile([]byte("drop table if exists users cascade; drop table if exists tokens cascade; drop table if exists remember_tokens;"), downFile)

	// Run Migrations
	err = doMigrate("up", "")
	if err != nil {
		exitGracefully(err)
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

	err = copyFileFromTemplate("templates/mailer/password-reset.html.tmpl", ade.RootPath+"/handlers/password-reset.html.tmpl")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFileFromTemplate("templates/mailer/password-reset.plain.tmpl", ade.RootPath+"/handlers/password-reset.plain.tmpl")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFileFromTemplate("templates/views/login.jet", ade.RootPath+"/views/views.jet")
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
