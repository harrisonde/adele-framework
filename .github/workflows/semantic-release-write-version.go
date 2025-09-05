package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

func main() {

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	file := path + "/adele.go"

	adele, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("Error" + err.Error())
	}

	body := string(adele)
	regx := regexp.MustCompile(`const Version = "(v.*)"`)
	raw := regx.ReplaceAllString(body, "const Version = \""+os.Getenv("GITHUB_REF_TAG")+"\"")

	err = ioutil.WriteFile(file, []byte(raw), 0644)
	if err != nil {
		fmt.Println("Error" + err.Error())
	}
}
