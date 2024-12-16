package main

import (
	"JSH-Automated-Stock-Trading/config"
	"JSH-Automated-Stock-Trading/service"
)

var ACCESS_TOKEN string

func main() {
	config.LoadConfigFromEnv()
	accessToken := service.GetAccessToken()
	service.GetBalancer(accessToken)
}
