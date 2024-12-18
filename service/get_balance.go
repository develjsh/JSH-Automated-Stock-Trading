package service

import (
	"JSH-Automated-Stock-Trading/config"
	"encoding/json"
	"fmt"
	"log"

	"github.com/valyala/fasthttp"
)

// 매수 가능 조회
func GetBalancer(accessToken string) {
	// 요청 헤더 설정
	headers := map[string]string{
		"Content-Type": "application/json",
		"appKey":       config.SetConfig.AppKey,
		"appSecret":    config.SetConfig.AppSecret,
	}

	// URL 설정
	PATH := "uapi/domestic-stock/v1/trading/inquire-psbl-order"
	url := fmt.Sprintf("%s/%s", config.SetConfig.UrlBase, PATH)

	// HTTP POST 요청 생성
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod("GET")
	req.SetRequestURI(url)
	req.Header.Set("Content-Type", headers["Content-Type"])
	req.Header.Set("appKey", headers["appKey"])
	req.Header.Set("appSecret", headers["appSecret"])
	req.Header.Set("authorization", "Bearer "+accessToken)
	req.Header.Set("tr_id", "VTTC8908R")
	req.Header.Set("custtype", "P")

	args := req.URI().QueryArgs()
	args.Set("CANO", config.SetConfig.Cano)               // 종합계좌번호
	args.Set("ACNT_PRDT_CD", config.SetConfig.AcntPrdtCd) // 계좌상품코드
	args.Set("PDNO", "005930")                            //상품번호
	args.Set("ORD_UNPR", "65500")                         // 주문단가
	args.Set("ORD_DVSN", "01")                            // 주문구분
	args.Set("CMA_EVLU_AMT_ICLD_YN", "Y")                 // CMA평가금액포함여부
	args.Set("OVRS_ICLD_YN", "Y")                         // 해외포함여부

	fmt.Printf("최종 요청 URI: %s\n", req.URI().String())

	// fasthttp 응답 객체 생성
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	// 요청 보내기
	if err := fasthttp.Do(req, resp); err != nil {
		log.Fatalf("Error sending request: %v", err)
	}

	// 응답 파싱
	var response map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		log.Fatalf("Error unmarshaling response: %v", err)
	}
	fmt.Println(response)
}
