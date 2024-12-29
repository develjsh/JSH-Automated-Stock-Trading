package service

import (
	"JSH-Automated-Stock-Trading/config"
	"JSH-Automated-Stock-Trading/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func SendMessage(msg string, webhookURL string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Send message error", r)
		}
	}()

	// 시간 msg 에 추가
	now := time.Now()
	formattedDate := now.Format("2006-01-02 15:04")
	msg = fmt.Sprintf("오늘 날짜: %s\n%s", formattedDate, msg)
	msg = utils.WrapWithSeparators(msg)
	// 메시지 생성
	message := map[string]string{
		"content": msg,
	}

	// 메시지를 JSON 형식으로 변환
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("Error marshaling message: %v", err)
	}

	// Discord 웹훅으로 POST 요청 보내기
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonMessage))
	if err != nil {
		log.Fatalf("Error sending message to Discord: %v", err)
	}
	defer resp.Body.Close()

	// 응답 상태 코드 출력
	fmt.Printf("Message sent: %s\n", message)
	fmt.Printf("Response status: %s\n", resp.Status)
}

func SendDiscordStartOfProgram(totalCash, targetBuyCount int, buyPercent float64) {
	// 종목당 매수 금액 계산
	buyAmount := float64(totalCash) * buyPercent
	// 메시지 생성
	message := fmt.Sprintf(
		"=== 국내 주식 자동매매 프로그램을 시작합니다 ===\n"+
			"==== 매수 전략 ====\n"+
			"총 보유 현금: %d원\n"+
			"종목당 매수 비율: %.2f%%\n"+
			"종목당 매수 금액: %.0f원\n"+
			"매수할 종목 수: %d개\n"+
			"=================\n",
		totalCash,
		buyPercent*100,
		buyAmount,
		targetBuyCount,
	)
	SendMessage(message, config.SetConfig.DiscordWebhookUrl)
}
