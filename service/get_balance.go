package service

import (
	"JSH-Automated-Stock-Trading/config"
	"encoding/json"
	"fmt"
	"log"

	"github.com/valyala/fasthttp"
)

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
	req.Header.Set("tr_id", "TTTC8908R")
	req.Header.Set("custtype", "P")

	args := req.URI().QueryArgs()
	args.Set("CANO", config.SetConfig.Cano)
	args.Set("ACNT_PRDT_CD", config.SetConfig.AcntPrdtCd)
	args.Set("PDNO", "005930")
	args.Set("ORD_UNPR", "65500")
	args.Set("ORD_DVSN", "01")
	args.Set("CMA_EVLU_AMT_ICLD_YN", "Y")
	args.Set("OVRS_ICLD_YN", "Y")

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
