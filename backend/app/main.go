package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"simple_account/app/Auths"
	"simple_account/app/Database"
	"simple_account/app/Email"
	Extension "simple_account/app/ExtensionChannel"
	"simple_account/app/Logger"
	Server "simple_account/app/MainServer"
	"syscall"
	"time"
)

type Config struct {
	Port                 int    `json:"listen_port"`
	LogsPath             string `json:"logs_dir_path"`
	ExtensionChannelPort int    `json:"extension_channel_port"`
	Database             struct {
		Source string `json:"source"`
	} `json:"database"`
	Email struct {
		Host                string `json:"email_host"`
		Port                int    `json:"port"`
		User                string `json:"user"`
		Password            string `json:"password"`
		ApiHost             string `json:"api_host"`
		FilesPath           string `json:"template_files_path"`
		VerificationDuraion int    `json:"verification_duration"`
	} `json:"email"`
	Auth struct {
		Total                 int    `json:"total_number_of_keys"`
		KeyFilesPath          string `json:"key_files_path"`
		ValidityDuration      int    `json:"token_validity_duration"`
		RenewTime             int    `json:"token_auto_renew_time"`
		CacheValidityDuration int    `json:"temporary_token_validity_duration"`
	} `json:"auth"`
}

func main() {

	configBytes, err := os.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	var config Config
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		panic(err)
	}

	logger, err := Logger.New(config.LogsPath)
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	database, err := Database.New(&Database.Config{
		SourceName: config.Database.Source,
	})
	if err != nil {
		panic(err)
	}

	email, err := Email.New(Email.Config{
		Host:      config.Email.Host,
		Port:      config.Email.Port,
		User:      config.Email.User,
		Password:  config.Email.Password,
		ApiHost:   config.Email.ApiHost,
		FilesPath: config.Email.FilesPath,
	})
	if err != nil {
		panic(err)
	}

	auth := &Auths.Auth{
		Total:                 config.Auth.Total,
		KeyFilesPath:          config.Auth.KeyFilesPath,
		ValidityDuration:      time.Duration(config.Auth.ValidityDuration) * 24 * time.Hour,
		RenewTime:             time.Duration(config.Auth.RenewTime) * 24 * time.Hour,
		CacheValidityDuration: time.Duration(config.Auth.CacheValidityDuration) * time.Hour,
	}

	err = auth.Init(logger)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server, err := Server.New(config.Port, database, email, auth, logger)
	if err != nil {
		panic(err)
	}

	if config.ExtensionChannelPort != 0 {
		extensionChannel, err := Extension.New(config.ExtensionChannelPort, database, email, auth, logger)
		if err != nil {
			panic(err)
		}
		go extensionChannel.Start()
		defer extensionChannel.Exit(ctx)
	}

	go server.Start()
	defer server.Exit(ctx)

	fmt.Println("START")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
