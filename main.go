package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
)

const VERSION = "4.0.0"

func main() {
	halt := DoOptions()
	if halt {
		return
	}
	log.SetFlags(log.Llongfile)
	if len(os.Args) < 2 {
		fmt.Println("missing argument: profile-alias")
	}
	alias := os.Args[1]

	//open and parse credentials
	file, err := ReadCredentialsFile()
	if err != nil {
		log.Fatalln(err)
	}
	credentials := Unmarshal(file)

	profile, err := GetCredentials(alias)
	if err == io.EOF {
		fmt.Println("no active session. Please login using aws sso login --profile {{profile}}")
		return
	}
	if err != nil {
		log.Println(err)
		return
	}
	//modify profile
	credentials[alias] = profile

	//re-write credentials file
	contents := credentials.Marshal()
	err = os.WriteFile(CredentialsFilepath(), contents, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("updated:", alias)
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
