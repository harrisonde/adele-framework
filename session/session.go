package session

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
)

type Session struct {
	CookieLifetime string
	CookiePersist  string
	CookieDomain   string
	CookieName     string
	SessionType    string
	DBPool         *sql.DB
	RedisPool      *redis.Pool
}

func (s *Session) InitSession() *scs.SessionManager {
	var persist, secure bool

	// sessions last for?
	minutes, err := strconv.Atoi(s.CookieLifetime)
	if err != nil {
		minutes = 60
	}

	// do cookies persist?
	if strings.ToLower(s.CookiePersist) == "true" {
		persist = true
	} else {
		persist = false
	}

	// cookies are secure?
	if strings.ToLower(s.CookiePersist) == "true" {
		secure = true
	}

	// create our session
	session := scs.New()
	session.Lifetime = time.Duration(minutes) * time.Minute
	session.Cookie.Persist = persist
	session.Cookie.Name = s.CookieName
	session.Cookie.Secure = secure
	session.Cookie.Domain = s.CookieDomain
	session.Cookie.SameSite = http.SameSiteLaxMode

	// What session store do we need to use
	switch strings.ToLower(s.SessionType) {
	case "redis":
		session.Store = redisstore.New(s.RedisPool)
	case "mysql", "mariadb":
		session.Store = mysqlstore.New(s.DBPool)
	case "postgres", "postgresql":
		session.Store = postgresstore.New(s.DBPool)
	default:
		// cookie
	}

	return session
}
