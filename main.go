package main

import (
	"JSH-Automated-Stock-Trading/config"
	"JSH-Automated-Stock-Trading/service"
	"fmt"
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

	// 매수 전략 설정
	symbolList := []string{"005930", "035720", "000660", "069500"}
	// 매수 완료된 종목 리스트
	var boughtList []string
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
		if t9.Before(tNow) && tNow.Before(tStart) && !soldout { // 잔여 수량 매도
			for sym, qty := range stockDict {
				if success := service.Sell(sym, qty, accessToken); success {
					fmt.Printf("매도 성공: %s, 수량: %d\n", sym, qty)
				} else {
					fmt.Printf("매도 실패: %s, 수량: %d\n", sym, qty)
				}
			}

			soldout = true
			boughtList = []string{}
			stockDict = service.GetStockBalance(accessToken)
		}
		if tStart.Before(tNow) && tNow.Before(tSell) { // AM 09:05 ~ PM 03:15 : 매수
			for _, sym := range symbolList {
				if len(boughtList) < targetBuyCount {
					// 이미 매수한 종목은 스킵
					isBought := false
					for _, bought := range boughtList {
						if bought == sym {
							isBought = true
							break
						}
					}
					if isBought {
						continue
					}
				}
			}
		}
		if tSell.Before(tNow) && tNow.Before(tExit) {
		}
		if tExit.Before(tNow) {
			service.SendMessage("프로그램을 종료합니다.", config.SetConfig.DiscordWebhookUrl)
			break
		}

	}

}
