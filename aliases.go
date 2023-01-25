package main

import (
	"log"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v3"
)

func GetAlias(name string) string {
	aliases := ReadConfigFile()
	var alias string
	if aliases == nil {
		log.Println("could not load profile aliases")
		alias = "default"
	} else {
		alias = aliases[name]
	}
	return alias
}

func ReadConfigFile() map[string]string {
	configFile := ConfigFilepath()

	if fileExists(configFile) {
		data, err := os.ReadFile(configFile)
		if err != nil {
			log.Fatalln("could not read config file: ", err)
		}
		type Config struct {
			Profiles map[string]string `yaml:"profiles"`
		}
		var config Config
		err = yaml.Unmarshal(data, &config)
		if err != nil {
			return nil
		}
		return config.Profiles
	}
	log.Printf("pcreds.yaml not found at %s\n", configFile)
	return nil
}

func ConfigFilepath() string {
	return filepath.Join(HomeDirectory(), ".aws", "credentials")
}
