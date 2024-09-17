package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	"log"
	"os"
	"strconv"
	"time"
)

func GetCredentials(alias string) Profile {
	accessToken, err := getAccessToken()
	if err != nil {
		log.Panic(err)
	}
	profileInfo := getProfile(alias)
	sess := sso.NewFromConfig(aws.Config{Region: "us-east-1"})
	output, err := sess.GetRoleCredentials(context.Background(), &sso.GetRoleCredentialsInput{
		AccessToken: &accessToken,
		AccountId:   &profileInfo.AccountId,
		RoleName:    &profileInfo.RoleName,
	})
	if err != nil {
		log.Panic(err)
	}
	profile := Profile{
		"aws_access_key_id":     *output.RoleCredentials.AccessKeyId,
		"aws_secret_access_key": *output.RoleCredentials.SecretAccessKey,
		"aws_session_token":     *output.RoleCredentials.SessionToken,
		"expiration":            strconv.FormatInt(output.RoleCredentials.Expiration, 10),
	}
	return profile
}

type CacheContents struct {
	AccessToken string `json:"accessToken"`
	ExpiresAt   string `json:"expiresAt"`
}

func getAccessToken() (string, error) {
	home, _ := os.UserHomeDir()
	cacheDir := home + "/.aws/sso/cache/"
	files, err := os.ReadDir(cacheDir)
	if err != nil {
		log.Println(err)
		return "", err
	}
	for _, file := range files {
		var contents CacheContents
		data, _ := os.ReadFile(cacheDir + file.Name())
		if err != nil {
			log.Println(err)
		}
		err := json.Unmarshal(data, &contents)
		if err != nil {
			log.Println(err)
		}
		if contents.AccessToken != "" && contents.ExpiresAt != "" {
			expiration, err := time.Parse(time.RFC3339, contents.ExpiresAt)
			if err != nil {
				log.Println(err)
			}
			if expiration.After(time.Now()) {
				return contents.AccessToken, nil
			} else {
				fmt.Println("expired")
			}
		}
	}
	return "", fmt.Errorf("could not find active session")
}

type ProfileInfo struct {
	AccountId string
	RoleName  string
}

func getProfile(alias string) ProfileInfo {
	home, _ := os.UserHomeDir()
	configPath := home + "/.aws/config"
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Println("cannot read config file")
		log.Panic(err)
	}
	profiles := Unmarshal(data)
	info := profiles["profile "+alias]
	return ProfileInfo{
		AccountId: info["sso_account_id"],
		RoleName:  info["sso_role_name"],
	}
}
