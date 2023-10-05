package adel

import "database/sql"

type initPaths struct {
	rootPath    string
	folderNames []string
}

type cookieConfig struct {
	name     string
	lifetime string
	persist  string
	secure   string
	domain   string
}

type databaseConfig struct {
	dsn      string
	database string
}

type Database struct {
	DataType string
	Pool     *sql.DB
}

type redisConfig struct {
	host     string
	password string
	prefix   string
}

type SubCommand struct {
	Name    string
	Help    string
	Hanlder interface{}
}

type Command struct {
	Name        string
	Help        string
	Description string
	Options     map[string]string
	Usage       string
}
