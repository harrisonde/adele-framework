package adel

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/CloudyKit/jet/v6"
	"github.com/alexedwards/scs/v2"
	"github.com/dgraph-io/badger/v3"
	"github.com/go-chi/chi/v5"
	"github.com/gomodule/redigo/redis"
	"github.com/harrisonde/adel/cache"
	"github.com/harrisonde/adel/filesystem/miniofilesystem"
	"github.com/harrisonde/adel/filesystem/s3filesystem"
	"github.com/harrisonde/adel/filesystem/sftpfilesystem"
	"github.com/harrisonde/adel/filesystem/webdavfilesystem"
	"github.com/harrisonde/adel/mailer"
	"github.com/harrisonde/adel/render"
	"github.com/harrisonde/adel/session"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

const verson = "1.0.0"

var myRedisCache *cache.RedisCache
var myBadgerCache *cache.BadgerCache
var redisPool *redis.Pool
var badgerPool *badger.DB

var sessionManager *scs.SessionManager

var maintenanceMode bool

type Adel struct {
	AppName       string
	Debug         bool
	Version       string
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
	RootPath      string
	Routes        *chi.Mux
	config        config // internal to the app, do not export
	Render        *render.Render
	JetViews      *jet.Set
	Session       *scs.SessionManager
	DB            Database
	EncryptionKey string
	Cache         cache.Cache
	Scheduler     *cron.Cron
	Mail          mailer.Mail
	Server        Server
	FileSystem    map[string]interface{}
	S3            s3filesystem.S3
	SFTP          sftpfilesystem.SFTP
	WebDAV        webdavfilesystem.WebDAV
	Minio         miniofilesystem.Minio
}

type Server struct {
	ServerName string
	Port       string
	Secure     bool
	URL        string
}

type config struct {
	port        string
	renderer    string
	cookie      cookieConfig
	sessionType string
	database    databaseConfig
	redis       redisConfig
	uploads     uploadConfig
}

type uploadConfig struct {
	allowedMimeTypes []string
	maxUploadSize    int64
}

// Called by project consuming our package
func (a *Adel) New(rootPath string) error {

	// Hold our root path and folder names
	pathConfig := initPaths{
		rootPath:    rootPath,
		folderNames: []string{"handlers", "migrations", "views", "mail", "data", "public", "tmp", "logs", "middleware", "screenshots"},
	}

	// Create folders
	err := a.Init(pathConfig)
	if err != nil {
		return err
	}

	// Check and read the environment from the .env
	err = a.checkDotEnv(rootPath)
	if err != nil {
		return err
	}

	err = godotenv.Load(rootPath + "/.env")
	if err != nil {
		return err
	}

	// Populate Adel values
	infoLog, errorLog := a.startLoggers()
	a.InfoLog = infoLog
	a.ErrorLog = errorLog
	a.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	a.Version = a.Version
	a.RootPath = rootPath
	a.Mail = a.createMailer()
	a.Routes = a.routes().(*chi.Mux) // Cast

	// file upload
	exploded := strings.Split(os.Getenv("FILE_TYPES_ALLOWED"), ",")
	var mimeTypes []string
	for _, m := range exploded {
		mimeTypes = append(mimeTypes, m)
	}

	var maxUploadSize int64
	if max, err := strconv.Atoi(os.Getenv("FILE_MAX_UPLOAD_SIZE")); err != nil {
		maxUploadSize = 10 << 20 // 10 mb
	} else {
		maxUploadSize = int64(max)
	}

	// Connect to database
	if os.Getenv("DATABASE_TYPE") != "" {
		db, err := a.OpenDB(os.Getenv("DATABASE_TYPE"), a.BuildDSN())
		if err != nil {
			errorLog.Println(err)
			os.Exit(1)
		}
		a.DB = Database{
			DataType: os.Getenv("DATABASE_TYPE"),
			Pool:     db,
		}
	}

	// Setup cron/scheduler
	scheduler := cron.New()
	a.Scheduler = scheduler

	// Connect to redis
	if os.Getenv("CACHE") == "redis" || os.Getenv("SESSION_TYPE") == "redis" {
		myRedisCache = a.createClientRedisCache()
		a.Cache = myRedisCache
		redisPool = myRedisCache.Conn
	}

	if os.Getenv("CACHE") == "badger" {
		myBadgerCache = a.createClientBadgerCache()
		a.Cache = myBadgerCache
		badgerPool = myBadgerCache.Conn

		_, err = a.Scheduler.AddFunc("@daily", func() {
			_ = myBadgerCache.Conn.RunValueLogGC(0.7)
		})
		if err != nil {
			return err
		}
	}

	// Application config
	a.config = config{
		port:     os.Getenv("PORT"),
		renderer: os.Getenv("RENDERER"),
		cookie: cookieConfig{
			name:     os.Getenv("COOKIE_NAME"),
			lifetime: os.Getenv("COOKIE_LIFETIME"),
			persist:  os.Getenv("COOKIE_PERSIST"),
			secure:   os.Getenv("COOKIE_SECURE"),
			domain:   os.Getenv("COOKIE_DOMAIN"),
		},
		sessionType: os.Getenv("SESSION_TYPE"),
		database: databaseConfig{
			database: os.Getenv("DATABASE_TYPE"),
			dsn:      a.BuildDSN(),
		},
		redis: redisConfig{
			host:     os.Getenv("REDIS_HOST"),
			password: os.Getenv("REDIS_PASSWORD"),
			prefix:   os.Getenv("REDIS_PREFIX"),
		},
		uploads: uploadConfig{
			maxUploadSize:    maxUploadSize,
			allowedMimeTypes: mimeTypes,
		},
	}

	secure := true
	if strings.ToLower(os.Getenv("SECURE")) == "false" {
		secure = false
	}

	a.Server = Server{
		ServerName: os.Getenv("SERVER_NAME"),
		Port:       os.Getenv("PORT"),
		Secure:     secure,
		URL:        os.Getenv("APP_URL"),
	}

	// Create session
	sess := session.Session{
		CookieLifetime: a.config.cookie.lifetime,
		CookiePersist:  a.config.cookie.persist,
		CookieName:     a.config.cookie.name,
		CookieDomain:   a.config.cookie.domain,
		SessionType:    a.config.sessionType,
	}

	switch a.config.sessionType {
	case "redis":
		sess.RedisPool = myRedisCache.Conn

	case "mysql", "postgres", "mariadb", "postgresql":
		sess.DBPool = a.DB.Pool
	}

	a.Session = sess.InitSession()
	a.EncryptionKey = os.Getenv("APP_KEY")

	// Create Jet engine
	if a.Debug {
		var views = jet.NewSet(
			jet.NewOSFileSystemLoader(fmt.Sprintf("%s/views", rootPath)),
			jet.InDevelopmentMode(),
		)
		a.JetViews = views
	} else {
		var views = jet.NewSet(
			jet.NewOSFileSystemLoader(fmt.Sprintf("%s/views", rootPath)),
		)
		a.JetViews = views
	}

	// Create renderer engine
	a.createRenderer()

	// Filesystem
	a.FileSystem = a.createFileSystem()

	// Start the mail channel
	go a.Mail.ListenForMail()

	return nil
}

// Create directories, if they do not already exist.
func (a *Adel) Init(p initPaths) error {
	root := p.rootPath
	for _, path := range p.folderNames {
		err := a.CreateDirIfNotExist(root + "/" + path)
		if err != nil {
			return err
		}
	}
	return nil
}

// Create env, if it does not already exist.
func (a *Adel) checkDotEnv(path string) error {
	err := a.CreateFileIfNotExist(fmt.Sprintf("%s", path))
	if err != nil {
		return err
	}
	return nil
}

// Create application logs
func (a *Adel) startLoggers() (*log.Logger, *log.Logger) {
	var infoLog *log.Logger
	var errorLog *log.Logger

	// Create
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	return infoLog, errorLog

}

func (a *Adel) createRenderer() {
	myRenderer := render.Render{
		Renderer: a.config.renderer,
		RootPath: a.RootPath,
		Port:     a.config.port,
		JetViews: a.JetViews,
		Session:  a.Session,
	}

	a.Render = &myRenderer
}

func (a *Adel) createMailer() mailer.Mail {
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	m := mailer.Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Templates:   a.RootPath + "/mail",
		Host:        os.Getenv("SMTP_HOST"),
		Port:        port,
		Username:    os.Getenv("SMTP_USERNAME"),
		Password:    os.Getenv("SMTP_PASSWORD"),
		Encryption:  os.Getenv("SMTP_ENCRYPTION"),
		FromName:    os.Getenv("FROM_NAME"),
		FromAddress: os.Getenv("FROM_ADDRESS"),
		Jobs:        make(chan mailer.Message, 20),
		Results:     make(chan mailer.Result, 20),
		API:         os.Getenv("MAILER_API"),
		APIKey:      os.Getenv("MAILER_KEY"),
		APIUrl:      os.Getenv("MAILER_URL"),
	}
	return m
}

func (a *Adel) BuildDSN() string {
	var dsn string

	switch os.Getenv("DATABASE_TYPE") {
	case "postgres", "postgresql":
		dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s timezone=UTC connect_timeout=5",
			os.Getenv("DATABASE_HOST"),
			os.Getenv("DATABASE_PORT"),
			os.Getenv("DATABASE_USER"),
			os.Getenv("DATABASE_NAME"),
			os.Getenv("DATABASE_SSL_MODE"))

		if os.Getenv("DATABASE_PASSWORD") != "" {
			dsn = fmt.Sprintf("%s password=%s", dsn, os.Getenv("DATABASE_PASSWORD"))
		}
	default:

	}
	return dsn
}

// Redis
func (a *Adel) createRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     50,
		MaxActive:   10000,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", a.config.redis.host,
				redis.DialPassword(a.config.redis.password))
		},

		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			_, err := conn.Do("PING")
			return err
		},
	}
}

// Create client redis cache
func (a *Adel) createClientRedisCache() *cache.RedisCache {
	cacheClient := cache.RedisCache{
		Conn:   a.createRedisPool(),
		Prefix: a.config.redis.prefix,
	}
	return &cacheClient
}

// Badger
func (a *Adel) createBadgerPool() *badger.DB {
	db, err := badger.Open(badger.DefaultOptions(a.RootPath + "/tmp/badger"))
	if err != nil {
		return nil
	}
	return db
}

// Create client for badger cache
func (a *Adel) createClientBadgerCache() *cache.BadgerCache {
	cacheClient := cache.BadgerCache{
		Conn: a.createBadgerPool(),
	}
	return &cacheClient
}

// Create the filesystem for the application
func (a *Adel) createFileSystem() map[string]interface{} {
	fileSystem := make(map[string]interface{})

	if os.Getenv("S3_KEY") != "" {
		s3 := s3filesystem.S3{
			Key:    os.Getenv("S3_KEY"),
			Secret: os.Getenv("S3_SECRET"),
			Region: os.Getenv("S3_REGION"),
			Bucket: os.Getenv("S3_BUCKET"),
		}
		fileSystem["S3"] = s3
		a.S3 = s3
	}

	if os.Getenv("MINIO_SECRET") != "" {
		useSSL := false
		if strings.ToLower(os.Getenv("MINIO_USESSL")) == "true" {
			useSSL = true
		}

		minio := miniofilesystem.Minio{
			Endpoint: os.Getenv("MINIO_ENDPOINT"),
			Key:      os.Getenv("MINIO_KEY"),
			Secret:   os.Getenv("MINIO_SECRET"),
			UseSSL:   useSSL,
			Region:   os.Getenv("MINIO_REGION"),
			Bucket:   os.Getenv("MINIO_BUCKET"),
		}
		fileSystem["MINIO"] = minio
		a.Minio = minio
	}

	if os.Getenv("SFTP_HOST") != "" {
		sftp := sftpfilesystem.SFTP{
			Host:     os.Getenv("SFTP_HOST"),
			User:     os.Getenv("SFTP_USER"),
			Password: os.Getenv("SFTP_PASSWORD"),
			Port:     os.Getenv("SFTP_PORT"),
		}
		fileSystem["SFTP"] = sftp
		a.SFTP = sftp
	}

	if os.Getenv("WEBDAV_HOST") != "" {
		webDAV := webdavfilesystem.WebDAV{
			Host:     os.Getenv("WEBDAV_HOST"),
			User:     os.Getenv("WEBDAV_USER"),
			Password: os.Getenv("WEBDAV_PASSWORD"),
		}
		fileSystem["WEBDAV"] = webDAV
		a.WebDAV = webDAV
	}

	return fileSystem

}

type RPCServer struct{}

func (r *RPCServer) MaintenanceMode(inMaintenanceMode bool, resp *string) error {
	if inMaintenanceMode {
		maintenanceMode = true
		*resp = "Server in maintenance mode"
	} else {
		maintenanceMode = false
		*resp = "Server live!"
	}
	return nil
}

func (a *Adel) listenRPC() {
	if os.Getenv("RPC_PORT") != "" {
		a.InfoLog.Println("Starting RPC server on port", os.Getenv("RPC_PORT"))
		err := rpc.Register(new(RPCServer))
		if err != nil {
			a.ErrorLog.Println(err)
			return
		}
		listen, err := net.Listen("tcp", "127.0.0.1:"+os.Getenv("RPC_PORT"))
		if err != nil {
			a.ErrorLog.Println(err)
			return
		}
		for {
			rpcConn, err := listen.Accept()
			if err != nil {
				continue
			}
			go rpc.ServeConn(rpcConn)

		}
	}
}
