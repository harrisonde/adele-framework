package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	adele, err := ioutil.ReadFile(path + "/adel.go")
	if err != nil {
		fmt.Println("Error" + err.Error())
	}

	body := string(adele)

	hasVersion := strings.ReplaceAll(body, "${{ADELE_RELEASE_VERSION}}", os.Getenv("GITHUB_REF_NAME"))

	err = ioutil.WriteFile(path+"/adel.go", []byte(hasVersion), 0644)
	if err != nil {
		fmt.Println("Error" + err.Error())
	}
}
