package main

import (
	"JSH-Automated-Stock-Trading/config"
	"JSH-Automated-Stock-Trading/service"
	"fmt"
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
	targetBuyCount := 3                          // 매수할 종목 수
	buyPercent := 0.33                           // 종목당 매수 금액 비율
	buyAmount := float64(totalCash) * buyPercent // 종목별 매수 금액 계산

	// 출력
	fmt.Println("==== 매수 전략 ====")
	fmt.Printf("총 보유 현금: %d원\n", totalCash)
	fmt.Printf("종목당 매수 비율: %.2f%%\n", buyPercent*100)
	fmt.Printf("종목당 매수 금액: %.0f원\n", buyAmount)
	fmt.Printf("매수할 종목 수: %d개\n", targetBuyCount)
	fmt.Println("=================")
}
