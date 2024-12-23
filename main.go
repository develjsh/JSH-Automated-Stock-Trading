package main

import (
	"JSH-Automated-Stock-Trading/config"
	"JSH-Automated-Stock-Trading/service"
	"time"
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
	// symbolList := []string{"005930", "035720", "000660", "069500"}
	// 매수 완료된 종목 리스트
	// var boughtList []string
	targetBuyCount := 3 // 매수할 종목 수
	buyPercent := 0.33  // 종목당 매수 금액 비율
	service.SendDiscordStartOfProgram(totalCash, targetBuyCount, buyPercent)
	soldout := false

	for {
		// 현재 시간
		tNow := time.Now()
		// 각 시간 설정
		t9 := time.Date(tNow.Year(), tNow.Month(), tNow.Day(), 9, 0, 0, 0, time.Local)
		tStart := time.Date(tNow.Year(), tNow.Month(), tNow.Day(), 9, 5, 0, 0, time.Local)
		tSell := time.Date(tNow.Year(), tNow.Month(), tNow.Day(), 15, 15, 0, 0, time.Local)
		tExit := time.Date(tNow.Year(), tNow.Month(), tNow.Day(), 15, 20, 0, 0, time.Local)

		// 오늘의 요일 (0=일요일, 6=토요일)
		today := int(tNow.Weekday())
		if today == 5 || today == 6 { // 토요일이나 일요일이면 자동 종료
			if err := service.SendMessage("주말이므로 프로그램을 종료합니다.", config.SetConfig.DiscordWebhookUrl); err != nil {
				break
			}
			break
		}
		if t9.Before(tNow) && tNow.Before(tStart) && !soldout {
		}
		if tStart.Before(tNow) && tNow.Before(tSell) {

		}
		if tSell.Before(tNow) && tNow.Before(tExit) {
		}
		if tExit.Before(tNow) {

		}

	}

}
