package service

import (
	"JSH-Automated-Stock-Trading/config"
	"encoding/json"
	"fmt"
	"log"

	"github.com/valyala/fasthttp"
)

func Sell(code, qty, accessToken string) bool {
	path := "uapi/domestic-stock/v1/trading/order-cash"
	url := fmt.Sprintf("%s/%s", config.SetConfig.UrlBase, path)

	// 요청 데이터
	data := map[string]interface{}{
		"CANO":         "your_account_number",     // 계좌번호
		"ACNT_PRDT_CD": "your_account_product_cd", // 계좌상품코드
		"PDNO":         code,                      // 종목 코드
		"ORD_DVSN":     "01",                      // 주문 구분 (시장가)
		"ORD_QTY":      qty,                       // 주문 수량
		"ORD_UNPR":     "0",                       // 주문 단가 (시장가)
	}

	// JSON 데이터 직렬화
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Error marshaling data: %v", err)
	}

	// fasthttp 요청 생성
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(url)
	req.Header.SetMethod("POST")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authorization", fmt.Sprintf("Bearer %s", "your_access_token"))
	req.Header.Set("appKey", "your_app_key")
	req.Header.Set("appSecret", "your_app_secret")
	req.Header.Set("tr_id", "TTTC0801U") // 모의 투자: VTTC0801U, 실전 투자: TTTC0801U
	req.Header.Set("custtype", "P")      // 개인: P, 법인: B
	req.Header.Set("hashkey", Hashkey(data))
	req.SetBody(jsonData)

	// fasthttp 응답 객체 생성
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	// 요청 보내기
	if err := fasthttp.Do(req, resp); err != nil {
		log.Fatalf("Error sending request: %v", err)
	}

	// 응답 처리
	var response map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		log.Fatalf("Error unmarshaling response: %v", err)
	}

	// 결과 처리
	if response["rt_cd"] == "0" {
		SendMessage(fmt.Sprintf("[매도 성공]%v", response), config.SetConfig.DiscordWebhookUrl)
		return true
	} else {
		SendMessage(fmt.Sprintf("[매도 실패]%v", response), config.SetConfig.DiscordWebhookUrl)
		return false
	}
}
