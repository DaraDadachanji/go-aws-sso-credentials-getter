package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
)

const VERSION = "4.0.0"

func main() {
	halt := DoOptions()
	if halt {
		return
	}

	paste := ReadClipboard()
	if len(paste) != 4 {
		log.Fatalln("expected 4 lines, received ", len(paste))
	}
	profileLine := paste[0]
	if !IsProfileName(profileLine) {
		log.Fatal("First line is not a profile tag")
	}
	name := ParseProfileName(profileLine)
	fmt.Println("received profile:", name)
	alias := GetAlias(name) //match against config file
	//open and parse credentials
	file, err := ReadCredentialsFile()
	if err != nil {
		log.Fatalln(err)
	}
	credentials := Unmarshal(file)

	//modify profile
	for _, line := range paste[1:] {
		key, value := ParseKeyValue(line)
		if credentials[alias] == nil {
			credentials[alias] = Profile{}
		}
		credentials[alias][key] = value
	}

	//re-write credentials file
	contents := credentials.Marshal()
	err = os.WriteFile(CredentialsFilepath(), contents, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("updated:", alias)
}

func ReadClipboard() []string {
	paste, _ := clipboard.ReadAll()
	reader := bufio.NewReader(strings.NewReader(paste))
	var lines []string
	for {
		var line []byte
		var readErr error
		isPrefix := true
		for isPrefix {
			var segment []byte
			segment, isPrefix, readErr = reader.ReadLine()
			line = append(line, segment...)
		}
		lines = append(lines, string(line))
		if readErr == io.EOF {
			return nonBlankLines(lines)
		}
	}
}

func nonBlankLines(lines []string) []string {
	var newLines []string
	for _, line := range lines {
		if !IsBlank(line) {
			newLines = append(newLines, line)
		}
	}
	return newLines
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
