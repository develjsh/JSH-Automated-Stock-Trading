package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func SendMessage(msg string, webhookURL string) {
	// 현재 시간 가져오기
	now := time.Now()

	// 메시지 생성
	message := map[string]string{
		"content": fmt.Sprintf("[%s] %s", now.Format("2006-01-02 15:04:05"), msg),
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
