package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"simple_account/app/Auths"
	"simple_account/app/Database"
	"simple_account/app/Email"
	Server "simple_account/app/Http"
	"simple_account/app/Logger"
	"syscall"
	"time"
)

type Config struct {
	Database struct {
		Source string `json:"source"`
	} `json:"database"`
	Email struct {
		Host     string `json:"email_host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		ApiHost  string `json:"api_host"`
	} `json:"email"`
	Auth struct {
		PrivateKeyFilePath    string `json:"private_key_file_path"`
		PublicKeyFilePath     string `json:"public_key_file_path"`
		ValidityDuration      int    `json:"token_validity_duration"`
		RenewTime             int    `json:"token_auto_renew_time"`
		CacheValidityDuration int    `json:"temporary_token_validity_duration"`
	} `json:"auth"`
}

func main() {

	Logger.Init()

	configBytes, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatalf("read json file fail: %v", err)
	}

	var config Config
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		panic(err)
	}

	database := Database.Config{
		SourceName: config.Database.Source,
	}

	email := Email.Config{
		Host:     config.Email.Host,
		Port:     config.Email.Port,
		User:     config.Email.User,
		Password: config.Email.Password,
		ApiHost:  config.Email.ApiHost,
	}

	auth := &Auths.Auth{
		PrivateKeyFilePath:    config.Auth.PrivateKeyFilePath,
		PublicKeyFilePath:     config.Auth.PublicKeyFilePath,
		ValidityDuration:      time.Duration(config.Auth.ValidityDuration) * 24 * time.Hour,
		RenewTime:             time.Duration(config.Auth.RenewTime) * 24 * time.Hour,
		CacheValidityDuration: time.Duration(config.Auth.CacheValidityDuration) * time.Hour,
	}
	err = auth.Init()
	if err != nil {
		panic(err)
	}

	server, err := Server.New(80, database, email, auth)

	if err != nil {
		panic(err)
	}

	go server.Start()

	fmt.Println("START")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	Logger.Close()

	if err := server.Exit(); err != nil {
		panic(err)
	}
}
