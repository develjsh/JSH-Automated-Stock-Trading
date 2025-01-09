package service

import (
	"JSH-Automated-Stock-Trading/config"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/valyala/fasthttp"
)

var response struct {
	Output struct {
		OrderableCash string `json:"ord_psbl_cash"`
	} `json:"output"`
}

// 주문 가능한 항목 조회
func GetBalancer(accessToken string) int {
	// URL 및 경로 설정
	path := "uapi/domestic-stock/v1/trading/inquire-psbl-order"
	url := fmt.Sprintf("%s/%s", config.SetConfig.UrlBase, path)

	// 요청 헤더 설정
	headers := map[string]string{
		"Content-Type":  "application/json",
		"authorization": "Bearer " + accessToken,
		"appKey":        config.SetConfig.AppKey,
		"appSecret":     config.SetConfig.AppSecret,
		"tr_id":         "VTTC8908R", // 실전 투자: TTTC8434R
		"custtype":      "P",
	}

	// 요청 파라미터 설정
	params := map[string]string{
		"CANO":                 config.SetConfig.Cano,       // 종합계좌번호
		"ACNT_PRDT_CD":         config.SetConfig.AcntPrdtCd, // 계좌상품코드
		"PDNO":                 "005930",                    // 상품번호
		"ORD_UNPR":             "65500",                     // 주문단가
		"ORD_DVSN":             "01",                        // 주문구분
		"CMA_EVLU_AMT_ICLD_YN": "Y",                         // CMA평가금액포함여부
		"OVRS_ICLD_YN":         "Y",                         // 해외포함여부
	}

	// HTTP GET 요청 생성
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod(fasthttp.MethodGet)
	req.SetRequestURI(url)

	// 헤더 설정
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// 쿼리 파라미터 설정
	query := req.URI().QueryArgs()
	for k, v := range params {
		query.Set(k, v)
	}

	// 응답 객체 생성
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	// 요청 전송
	if err := fasthttp.Do(req, resp); err != nil {
		SendMessage("InquireBalance 요청 전송 실패", config.SetConfig.DiscordWebhookUrl)
		return -1
	}

	// JSON 파싱
	if err := json.Unmarshal(resp.Body(), &getBalancerResponse); err != nil {
		log.Fatalf("Error unmarshaling response: %v", err)
	}

	// 주문 가능 현금 잔고 가져오기
	orderableCash := getBalancerResponse.Output.OrderableCash
	fmt.Printf("주문 가능 현금 잔고: %s원\n", orderableCash)

	// 문자열을 정수로 변환
	cash, err := strconv.Atoi(orderableCash)
	if err != nil {
		log.Fatalf("Error converting orderable cash to int: %v", err)
	}

	return cash
}
