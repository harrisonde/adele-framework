package adele

import (
	"log"

	"github.com/gobuffalo/pop"
	_ "github.com/gobuffalo/pop"

	"github.com/golang-migrate/migrate/v4"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func (a *Adele) PopConnect() (*pop.Connection, error) {
	tx, err := pop.Connect("development") // TODO: Do we want to default to development? Seems to me that a env pivot is helpful.
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (a *Adele) CreatePopMigration(up, down []byte, migrationName, migrationType string) error {
	var migrationPath = a.RootPath + "/migrations"
	err := pop.MigrationCreate(migrationPath, migrationName, migrationType, up, down)
	if err != nil {
		return err
	}
	return nil
}

func (a *Adele) RunPopMigrations(tx *pop.Connection) error {
	var migrationPath = a.RootPath + "/migrations"

	fm, err := pop.NewFileMigrator(migrationPath, tx)
	if err != nil {
		return err
	}

	err = fm.Up()
	if err != nil {
		return err
	}

	return nil
}

func (a *Adele) PopMigrateDown(tx *pop.Connection, steps ...int) error {
	var migrationPath = a.RootPath + "/migrations"

	step := 1
	if len(steps) > 0 {
		step = steps[0]
	}

	fm, err := pop.NewFileMigrator(migrationPath, tx)
	if err != nil {
		return err
	}

	err = fm.Down(step)
	if err != nil {
		return err
	}
	return nil
}

func (a *Adele) PopMigrateReset(tx *pop.Connection) error {
	var migrationPath = a.RootPath + "/migrations"
	fm, err := pop.NewFileMigrator(migrationPath, tx)
	if err != nil {
		return err
	}
	err = fm.Reset()
	if err != nil {
		return err
	}
	return nil
}

func (c *Adele) MigrateUp(dsn string) error {
	m, err := migrate.New("file://"+c.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil {
		log.Println("Error running migration: ", err)
		return err
	}
	return nil
}

func (c *Adele) MigrateDownAll(dsn string) error {
	m, err := migrate.New("file://"+c.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Down(); err != nil {
		return err
	}
	return nil
}

func (c *Adele) Steps(n int, dsn string) error {
	m, err := migrate.New("file://"+c.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Steps(n); err != nil {
		return err
	}
	return nil
}

func (c *Adele) MigrateForce(dsn string) error {
	m, err := migrate.New("file://"+c.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Force(-1); err != nil {
		return err
	}

	return nil
}
