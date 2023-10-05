package main

import (
	"errors"
	"strings"

	"github.com/harrisonde/adel"
)

var MakeMailCommand = &adel.Command{
	Name:        "make mail",
	Help:        "create a new mail template",
	Description: "use the make mail command to create a new mail template in the mail directory",
	Usage:       "make mail <name>",
	Options:     map[string]string{},
}

func doMakeMail(arg3 string) error {
	if arg3 == "" {
		return errors.New("you must provide a name for the mail template")
	}

	htmlMail := ade.RootPath + "/mail/" + strings.ToLower(arg3) + ".html.tmpl"
	plainMail := ade.RootPath + "/mail/" + strings.ToLower(arg3) + ".plain.tmpl"

	err := copyFileFromTemplate("templates/mailer/mail.html.tmpl", htmlMail)
	if err != nil {
		return err
	}

	err = copyFileFromTemplate("templates/mailer/mail.plain.tmpl", plainMail)
	if err != nil {
		return err
	}
	return nil
}
