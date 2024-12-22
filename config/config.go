package config

import (
	"fmt"
	"log"
	"os"
	"time"
)

var SetConfig Config

type Config struct {
	Location          string
	AppKey            string
	AppSecret         string
	Cano              string
	AcntPrdtCd        string
	UrlBase           string
	DiscordWebhookUrl string
}

func LoadConfigFromEnv() {
	SetConfig = Config{
		Location:          os.Getenv("LOCATION"),
		AppKey:            os.Getenv("APP_KEY"),
		AppSecret:         os.Getenv("APP_SECRET"),
		Cano:              os.Getenv("CANO"),
		AcntPrdtCd:        os.Getenv("ACNT_PRDT_CD"),
		UrlBase:           os.Getenv("URL_BASE"),
		DiscordWebhookUrl: os.Getenv("DISCORD_WEBHOOK_URL"),
	}

	// 시간 설정
	setLocationTime(SetConfig.Location)

	// 필수 환경 변수가 설정되지 않았는지 확인
	if SetConfig.AppKey == "" || SetConfig.AppSecret == "" {
		log.Fatal("환경 변수 설정이 올바르지 않습니다. 필수 변수를 확인하세요.")
	}
}

func setLocationTime(location string) {
	if location == "" {
		location = "Asia/Seoul"
	}

	loc, err := time.LoadLocation(location)
	if err != nil {
		fmt.Printf("Failed to load location: %v, defaulting to Asia/Seoul\n", err)
		loc, _ := time.LoadLocation("Asia/Seoul")
		time.Local = loc
	}
	time.Local = loc
}
