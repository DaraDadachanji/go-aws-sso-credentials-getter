package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"
)

type Profiles map[string]Profile

type Profile map[string]string

func (p *Profiles) Marshal() []byte {
	var lines []string
	for name, profile := range *p {
		line := fmt.Sprintf("[%s]", name)
		lines = append(lines, line)
		for key, value := range profile {
			line = fmt.Sprintf("%s = %s", key, value)
			lines = append(lines, line)
		}
	}
	contents := []byte(strings.Join(lines, "\n"))
	contents = append(contents, '\n') //trailing newline
	return contents

}

func Unmarshal(data []byte) Profiles {
	reader := bufio.NewReader(bytes.NewReader(data))
	profiles := Profiles{}
	for {
		var profileName string
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if IsBlank(line) {
			continue
		}
		if IsProfileName(line) {
			profileName = ParseProfileName(line)
			profiles[profileName] = Profile{}
			continue
		}
		key, value := ParseKeyValue(line)
		profiles[profileName][key] = value
	}
	return Profiles{}
}

func IsBlank(line string) bool {
	match, _ := regexp.MatchString(`^\s*$`, line)
	return match
}

func IsProfileName(line string) bool {
	match, _ := regexp.MatchString(`^\[[A-Za-z0-9\-_]+\]`, line)
	return match
}

func ParseProfileName(line string) string {
	profile := line[1 : len(line)-2]
	return profile
}

func ParseKeyValue(line string) (key string, value string) {
	r := regexp.MustCompile(`(?P<key>[^= ]*)[ ]*=[ ]*(?P<value>"[^" ]*"|[^," ]*)`)
	parts := r.FindStringSubmatch(line)
	if len(parts) < 3 {
		log.Panic("could not parse credentials file")
	}
	return parts[1], parts[2]
}
