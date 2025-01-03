package main

import (
	"JSH-Automated-Stock-Trading/config"
	"JSH-Automated-Stock-Trading/service"
	"fmt"
	"time"
)

var ACCESS_TOKEN string

func main() {
	defer func() {
		if r := recover(); r != nil {
			service.SendMessage("프로그램 실행 중 panic 발생!!!", config.SetConfig.DiscordWebhookUrl)
		}
	}()

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
	targetBuyCount := 3                               // 매수할 종목 수
	buyPercent := 0.33                                // 종목당 매수 금액 비율
	buyAmount := int(float64(totalCash) * buyPercent) // 종목별 주문 금액 계산
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
			service.SendMessage("주말이므로 프로그램을 종료합니다.", config.SetConfig.DiscordWebhookUrl)
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
					targetPrice, err := service.GetTargetPrice(sym, accessToken)
					if err != nil {
						continue
					}
					currentPrice, err := service.GetCurrentPrice(sym, accessToken)
					if err != nil {
						continue
					}
					if targetPrice < currentPrice {
						buyQty := 0 // 매수할 수량 초기화
						buyQty = int(float64(buyAmount) / float64(currentPrice))
						if buyQty > 0 {
							message := fmt.Sprintf("%s 목표가 달성(%d < %d) 매수를 시도합니다.", sym, targetPrice, currentPrice)
							service.SendMessage(message, config.SetConfig.DiscordWebhookUrl)

							// 매수 시도
							if success := service.Buy(sym, buyQty, accessToken); success {
								soldout = false
								boughtList = append(boughtList, sym)
								stockDict = service.GetStockBalance(accessToken)
							}
						}
					}
					time.Sleep(1 * time.Second)
				}
			}
			time.Sleep(1 * time.Second)
			// 매 30분마다 잔고 확인
			if tNow.Minute()%30 == 0 && tNow.Second() <= 5 {
				stockDict = service.GetStockBalance(accessToken)
				time.Sleep(5 * time.Second)
			}
		}
		if tSell.Before(tNow) && tNow.Before(tExit) { // PM 03:15 ~ PM 03:20 : 일괄 매도
			if !soldout {
				stockDict = service.GetStockBalance(accessToken)
				stockDictSize := len(stockDict)
				sellSize := 0
				for sym, qty := range stockDict {
					if success := service.Sell(sym, qty, accessToken); success {
						sellSize += 1
						service.SendMessage(fmt.Sprintf("매도 성공: %s, 수량: %s\n", sym, qty), config.SetConfig.DiscordWebhookUrl)
					} else {
						service.SendMessage(fmt.Sprintf("매도 실패: %s, 수량: %s\n", sym, qty), config.SetConfig.DiscordWebhookUrl)
					}
				}
				if stockDictSize == sellSize {
					soldout = true
				}
			}
		}
		if tExit.Before(tNow) { // PM 03:20 ~ :프로그램 종료
			service.SendMessage("프로그램을 종료합니다.", config.SetConfig.DiscordWebhookUrl)
			break
		}
	}

}
