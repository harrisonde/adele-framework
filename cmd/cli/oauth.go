package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/harrisonde/adele"
)

var OauthCommand = &adele.Command{
	Name:        "make oauth",
	Help:        "install oauth",
	Description: "install oauth2 authentication into your adele application",
	Usage:       "make oauth",
	Options:     map[string]string{},
}

func doOauth() error {
	fmt.Printf("Adele oAuth")

	checkForDb()

	color.Yellow("\n\nStarting installation")

	dbType := ade.DB.DataType

	tx, err := ade.PopConnect()
	if err != nil {
		exitGracefully(err)
	}
	defer tx.Close()

	upBytes, err := templateFS.ReadFile("templates/migrations/oauth_tables." + dbType + ".sql")
	if err != nil {
		exitGracefully(err)
	}

	downBytes := []byte("drop table if exists tokens cascade;")
	if err != nil {
		exitGracefully(err)
	}

	color.Green("  Creating migrations...")
	err = ade.CreatePopMigration(upBytes, downBytes, "oauth", "sql")
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
		"middleware",
		"migrations",
	}

	root := ade.RootPath
	for _, path := range appDirs {
		err := ade.CreateDirIfNotExist(root + "/" + path)
		if err != nil {
			return err
		}
	}

	color.Green("  Creating models...")

	err = copyFileFromTemplate("templates/data/token.go.txt", ade.RootPath+"/data/token.go")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFileFromTemplate("templates/data/client.go.txt", ade.RootPath+"/data/client.go")
	if err != nil {
		exitGracefully(err)
	}

	color.Green("  Creating middleware...")

	err = copyFileFromTemplate("templates/middleware/auth-token.go.txt", ade.RootPath+"/middleware/auth-token.go")
	if err != nil {
		exitGracefully(err)
	}

	color.Green("  Creating handlers...")

	err = copyFileFromTemplate("templates/handlers/oauth-handlers.go.txt", ade.RootPath+"/handlers/oauth-handlers.go")
	if err != nil {
		exitGracefully(err)
	}

	color.Green("  Creating commands...")

	data, err := templateFS.ReadFile("templates/cmd/client.go.txt")
	if err != nil {
		return err
	}
	fileName := ade.RootPath + "/cmd/client.go"
	handler := string(data)
	handler = strings.ReplaceAll(handler, "$APPNAME$", ade.AppName)
	err = ioutil.WriteFile(fileName, []byte(handler), 0644)
	if err != nil {
		return err
	}

	cmd := exec.Command("go", "get", "github.com/harrisonde/adele")
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
	color.Green("  1. add client and token models to your application:")
	fmt.Printf("    models.go")
	fmt.Printf("\n      type Models struct {")
	fmt.Printf("\n        Users  User")
	fmt.Printf("\n        ... ")
	fmt.Printf("\n        Clients Client")
	fmt.Printf("\n        Tokens  Token")
	fmt.Printf("\n      }")
	fmt.Printf("\n      return Models{")
	fmt.Printf("\n        Users:    User{},")
	fmt.Printf("\n        ... ")
	fmt.Printf("\n        Clients:  Client{},")
	fmt.Printf("\n        Tokens:   Token{},")
	fmt.Printf("\n      }")
	color.Green("\n  2. add client command to your commands:")
	fmt.Printf("    commands.go")
	fmt.Printf("\n      case \"oauth\":")
	fmt.Printf("\n        if arg2 == \"client\" {")
	fmt.Printf("\n          return c.doCreateOauthClient(arg1, arg2, arg3)")
	fmt.Printf("\n        }")
	color.Green("\n  3. add token to user struct:")
	fmt.Printf("    users.go")
	fmt.Printf("\n      type User struct {")
	fmt.Printf("\n        ID        int       `db:\"id,omitempty\"`")
	fmt.Printf("\n        FirstName string    `db:\"first_name\"`")
	fmt.Printf("\n        ... ")
	fmt.Printf("\n        Token     Token     `db:\"-\"`")
	color.Green("\n  4. add oauth exchange route to your application:")
	fmt.Printf("    routes-web.go")
	fmt.Printf("\n      routes:")
	fmt.Printf("\n        r.Post(\"/oauth/token\", a.Handlers.PasswordGrantExchange)")
	color.Green("\n  5. create your first client")
	fmt.Printf("\n    $ adele oauth client <client name> --user=<id>")
	fmt.Printf("\n")

	return nil
}
