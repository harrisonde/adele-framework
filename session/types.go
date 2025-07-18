package session

import (
	"database/sql"

	"github.com/gomodule/redigo/redis"
)

type Session struct {
	CookieLifetime string
	CookiePersist  string
	CookieDomain   string
	CookieSecure   string
	CookieName     string
	DBPool         *sql.DB
	RedisPool      *redis.Pool
	SessionType    string
}
