package service

import (
	"JSH-Automated-Stock-Trading/config"
	"encoding/json"
	"fmt"
	"log"

	"github.com/valyala/fasthttp"
)

func hashkey(datas map[string]interface{}) string {
	// 요청 헤더 설정
	headers := map[string]string{
		"Content-Type": "application/json",
		"appKey":       config.SetConfig.AppKey,
		"appSecret":    config.SetConfig.AppSecret,
	}

	// 요청 본문을 JSON으로 변환
	bodyBytes, err := json.Marshal(datas)
	if err != nil {
		log.Fatalf("Error marshaling data: %v", err)
	}

	// URL 설정
	url := fmt.Sprintf("%s/uapi/hashkey", config.SetConfig.UrlBase)

	// HTTP POST 요청 생성
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.Set("Content-Type", headers["Content-Type"])
	req.Header.Set("appKey", headers["appKey"])
	req.Header.Set("appSecret", headers["appSecret"])
	req.SetRequestURI(url)
	req.Header.SetMethod("POST")
	req.SetBody(bodyBytes)

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

	// "HASH" 값 추출
	if hash, ok := response["HASH"].(string); ok {
		return hash
	}
	return ""
}
