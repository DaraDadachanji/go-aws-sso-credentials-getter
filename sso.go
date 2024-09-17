package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func GetCredentials(alias string) (Profile, error) {
	accessToken, err := getAccessToken()
	if err != nil {
		return nil, err
	}
	profileInfo := getProfile(alias)
	sess := sso.NewFromConfig(aws.Config{Region: "us-east-1"})
	output, err := sess.GetRoleCredentials(context.Background(), &sso.GetRoleCredentialsInput{
		AccessToken: &accessToken,
		AccountId:   &profileInfo.AccountId,
		RoleName:    &profileInfo.RoleName,
	})
	if err != nil {
		return nil, err
	}
	profile := Profile{
		"aws_access_key_id":     *output.RoleCredentials.AccessKeyId,
		"aws_secret_access_key": *output.RoleCredentials.SecretAccessKey,
		"aws_session_token":     *output.RoleCredentials.SessionToken,
		"expiration":            strconv.FormatInt(output.RoleCredentials.Expiration, 10),
	}
	return profile, nil
}

type CacheContents struct {
	AccessToken string `json:"accessToken"`
	ExpiresAt   string `json:"expiresAt"`
}

func getAccessToken() (string, error) {
	cacheDir := filepath.Join(HomeDirectory(), ".aws", "sso", "cache")
	files, err := os.ReadDir(cacheDir)
	if err != nil {
		log.Println(err)
		return "", err
	}
	for _, file := range files {
		var contents CacheContents
		data, _ := os.ReadFile(filepath.Join(cacheDir, file.Name()))
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
			}
		}
	}
	fmt.Println("could not find active session")
	return "", io.EOF
}

type ProfileInfo struct {
	AccountId string
	RoleName  string
}

func getProfile(alias string) ProfileInfo {
	configPath := filepath.Join(HomeDirectory(), ".aws", "config")
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
