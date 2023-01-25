package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
	"gopkg.in/ini.v1"
)

const VERSION = "3.0.0"

func main() {
	halt := DoOptions()
	if halt {
		return
	}

	paste, _ := clipboard.ReadAll()
	incoming, err := ini.Load([]byte(paste))
	if err != nil {
		log.Println(err)
		log.Fatalln("could not parse clipboard")
	}
	if len(incoming.Sections()) != 2 {
		log.Fatalln("invalid credentials format")
	}
	section := incoming.Sections()[1] // 0th section is DEFAULT
	name := section.Name()
	log.Println("recieved:", name)
	alias := GetAlias(name) //match against config file
	//open and parse credentials
	file, err := ReadCredentialsFile()
	if err != nil {
		log.Fatalln(err)
	}
	credentials, err := ini.Load(file)
	if err != nil {
		log.Println(err)
		log.Fatalln("could not parse credentials file")
	}

	//modify profile
	if credentials.HasSection(alias) {
		credentials.DeleteSection(alias)
	}
	newSection, err := credentials.NewSection(alias)
	if err != nil {
		log.Fatalln(err)
	}
	newSection.SetBody(section.Body())

	//re-write credentials file
	var buff bytes.Buffer
	_, err = credentials.WriteTo(&buff)
	if err != nil {
		log.Fatalln(err)
	}
	contents, err := ioutil.ReadAll(&buff)
	if err != nil {
		log.Fatalln(err)
	}
	err = os.WriteFile(CredentialsFilepath(), contents, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("updated:", alias)
}

func ReadClipboard() []string {
	paste, _ := clipboard.ReadAll()
	reader := bufio.NewReader(strings.NewReader(paste))
	var lines []string
	for {
		line, readErr := reader.ReadString('\n')
		lines = append(lines, line)
		if readErr == io.EOF {
			return lines
		}
	}
}

func ReadCredentialsFile() ([]byte, error) {
	credentialsFile := CredentialsFilepath()
	if !fileExists(credentialsFile) {
		return nil, fmt.Errorf("file not found: %s", credentialsFile)
	}
	data, err := os.ReadFile(credentialsFile)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) || info.IsDir() {
		return false
	}
	return true
}

func HomeDirectory() string {
	u, _ := user.Current()
	return u.HomeDir
}

func CredentialsFilepath() string {
	return filepath.Join(HomeDirectory(), ".aws", "/credentials")
}
