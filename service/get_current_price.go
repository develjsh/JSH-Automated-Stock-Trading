package service

import (
	"JSH-Automated-Stock-Trading/config"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/valyala/fasthttp"
)

func GetCurrentPrice(code string, accessToken string) (int, error) {
	// URL 및 경로 설정
	path := "uapi/domestic-stock/v1/quotations/inquire-price"
	url := fmt.Sprintf("%s/%s", config.SetConfig.UrlBase, path)

	// 요청 헤더 설정
	headers := map[string]string{
		"Content-Type":  "application/json",
		"authorization": "Bearer " + accessToken,
		"appKey":        config.SetConfig.AppKey,
		"appSecret":     config.SetConfig.AppSecret,
		"tr_id":         "FHKST01010100",
	}

	// 요청 파라미터 설정
	params := map[string]string{
		"fid_cond_mrkt_div_code": "J",
		"fid_input_iscd":         code,
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
		return 0, fmt.Errorf("error sending request: %v", err)
	}

	// 응답 데이터 파싱
	var result struct {
		Output struct {
			StckPrpr string `json:"stck_prpr"` // 현재가
		} `json:"output"`
	}

	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return 0, fmt.Errorf("error unmarshaling response: %v", err)
	}

	// 현재가 반환
	currentPrice, err := strconv.Atoi(result.Output.StckPrpr)
	if err != nil {
		return 0, fmt.Errorf("error converting stck_prpr to int: %v", err)
	}

	return currentPrice, nil
}
