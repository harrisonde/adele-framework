package main

import (
	"fmt"
	"time"

	"github.com/harrisonde/adel"
)

var MakeSessionCommand = &adel.Command{
	Name: "make session",
	Help: "create a table in the database to store sessions",
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

	err = doMigrate("up", "")
	if err != nil {
		exitGracefully(err)
	}

	return nil
}
