package mailer

import (
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"log"
	"os"
	"testing"
	"time"
)

var pool *dockertest.Pool
var resource *dockertest.Resource

var mailer = Mail{
	Domain:      "localhost",
	Templates:   "./testdata/mail",
	Host:        "localhost",
	Port:        1026,
	Encryption:  "none",
	FromAddress: "me@here.com",
	FromName:    "Adel",
	Jobs:        make(chan Message, 1),
	Results:     make(chan Result, 1),
}

func TestMain(m *testing.M) {
	p, err := dockertest.NewPool("")
	if err != nil {
		log.Fatal("can not connect to Docker", err)
	}

	pool = p

	options := dockertest.RunOptions{
		Repository:   "mailhog/mailhog",
		Tag:          "latest",
		Env:          []string{},
		ExposedPorts: []string{"1025", "8025"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"1025": {
				{HostIP: "0.0.0.0", HostPort: "1026"},
			},
			"8025": {
				{HostIP: "0.0.0.0", HostPort: "8026"},
			},
		},
	}

	resource, err := pool.RunWithOptions(&options)
	if err != nil {
		log.Println(err)
		_ = pool.Purge(resource)
		log.Fatal("Could not start resource")
	}

	// Cant ping Mailhog, so we sleep.
	// Any way we can side-step a sleep?

	time.Sleep(2 * time.Second)

	go mailer.ListenForMail()

	code := m.Run()

	// Kill docker image
	if err := pool.Purge(resource); err != nil {
		log.Fatal("could not purge resource", err)
	}

	os.Exit(code)
}
