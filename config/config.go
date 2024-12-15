package config

import (
	"log"
	"os"
)

type Config struct {
	AppKey            string
	AppSecret         string
	Cano              string
	AcntPrdtCd        string
	UrlBase           string
	DiscordWebhookUrl string
}

func LoadConfigFromEnv() Config {
	config := Config{
		AppKey:            os.Getenv("APP_KEY"),
		AppSecret:         os.Getenv("APP_SECRET"),
		Cano:              os.Getenv("CANO"),
		AcntPrdtCd:        os.Getenv("ACNT_PRDT_CD"),
		UrlBase:           os.Getenv("URL_BASE"),
		DiscordWebhookUrl: os.Getenv("DISCORD_WEBHOOK_URL"),
	}

	// 필수 환경 변수가 설정되지 않았는지 확인
	if config.AppKey == "" || config.AppSecret == "" {
		log.Fatal("환경 변수 설정이 올바르지 않습니다. 필수 변수를 확인하세요.")
	}

	return config
}
