package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v3"
)

func GetAlias(name string) string {
	config := ReadConfigFile()
	alias, ok := config.Profiles[name]
	if !ok {
		fmt.Printf("could not find alias for profile: %s\n", name)
		alias = config.PromptForAlias(name)
	}
	return alias
}

func (c *Config) PromptForAlias(name string) string {
	var userInput string
	for {
		fmt.Print("Please choose an alias for this profile: ")
		fmt.Scanln(&userInput)
		if len(userInput) != 0 {
			return userInput
		} else {
			log.Println("name must be greater than 0 characters")
		}
	}
}

func (c *Config) Write() {
	data, _ := yaml.Marshal(c)
	err := os.WriteFile(ConfigFilepath(), data, 0644)
	if err != nil {
		log.Println("failed to write new config file", err)
	}
}

type Config struct {
	Profiles map[string]string `yaml:"profiles"`
}

func ReadConfigFile() Config {
	configFile := ConfigFilepath()

	if fileExists(configFile) {
		data, err := os.ReadFile(configFile)
		if err != nil {
			log.Fatalln("could not read config file: ", err)
		}

		var config Config
		err = yaml.Unmarshal(data, &config)
		if err != nil {
			return Config{Profiles: map[string]string{}}
		}
		return config
	} else {
		log.Printf("pcreds.yaml not found at %s\n", configFile)
		return Config{Profiles: map[string]string{}}
	}
}

func ConfigFilepath() string {
	return filepath.Join(HomeDirectory(), ".aws", "credentials")
}
