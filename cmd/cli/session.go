package main

import (
	"fmt"
	"time"

	"github.com/harrisonde/adele-framework"
)

var MakeSessionCommand = &adele.Command{
	Name:        "make session",
	Help:        "install sessions",
	Description: "use the make session command to install session into your application; creates a table in the database to store sessions and runs migrations",
	Usage:       "make session",
	Options: map[string]string{
		"-s, --skip": "skip running up migration",
	},
}

func doSessionTable() error {
	dbType := ade.DB.DataType

	if dbType == "mariadb" {
		dbType = "mysql"
	}

	if dbType == "postgresql" {
		dbType = "postgres"
	}

	fileName := fmt.Sprintf("%d_create_session_table", time.Now().UnixMicro())

	upFile := ade.RootPath + "/migrations/" + fileName + "." + dbType + ".up.sql"
	downFile := ade.RootPath + "/migrations/" + fileName + "." + dbType + ".down.sql"

	err := copyFileFromTemplate("templates/migrations/"+dbType+"_session.sql", upFile)
	if err != nil {
		exitGracefully(err)
	}

	err = copyDataToFile([]byte("drop table sessions"), downFile)
	if err != nil {
		exitGracefully(err)
	}

	longOption, _ := GetOption("skip")
	shortOption, _ := GetOption("s")
	if longOption == "skip" || shortOption == "s" {
		return nil
	}

	err = doMigrate("up", "")
	if err != nil {
		exitGracefully(err)
	}

	return nil
}
