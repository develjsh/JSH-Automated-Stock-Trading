package main

import (
	"JSH-Automated-Stock-Trading/config"
	"fmt"
)

var ACCESS_TOKEN string

func main() {
	config := config.LoadConfigFromEnv()
	fmt.Println(config)
}
