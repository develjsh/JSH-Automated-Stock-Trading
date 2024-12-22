package main

import (
	"JSH-Automated-Stock-Trading/config"
	"JSH-Automated-Stock-Trading/service"
)

var ACCESS_TOKEN string

func main() {
	config.LoadConfigFromEnv()
	accessToken := service.GetAccessToken()

	// 보유 현금 조회
	totalCash := service.GetBalancer(accessToken)

	// 주식 잔고 조회
	stockDict := service.GetStockBalance(accessToken)

	// 보유 주식 목록 생성
	var boughtList []string
	for sym := range stockDict {
		boughtList = append(boughtList, sym)
	}

	// 매수 전략 설정
	targetBuyCount := 3 // 매수할 종목 수
	buyPercent := 0.33  // 종목당 매수 금액 비율
	service.SendDiscordStartOfProgram(totalCash, targetBuyCount, buyPercent)

}
