package service

import (
	"JSH-Automated-Stock-Trading/config"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/valyala/fasthttp"
)

func GetTargetPrice(code string, accessToken string) (float64, error) {
	// URL 및 경로 설정
	path := "uapi/domestic-stock/v1/quotations/inquire-daily-price"
	url := fmt.Sprintf("%s/%s", config.SetConfig.UrlBase, path)

	// 요청 헤더 설정
	headers := map[string]string{
		"Content-Type":  "application/json",
		"authorization": "Bearer " + accessToken,
		"appKey":        config.SetConfig.AppKey,
		"appSecret":     config.SetConfig.AppSecret,
		"tr_id":         "FHKST01010400",
	}

	// 요청 파라미터 설정
	params := map[string]string{
		"fid_cond_mrkt_div_code": "J",
		"fid_input_iscd":         code,
		"fid_org_adj_prc":        "1",
		"fid_period_div_code":    "D",
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
		Output []struct {
			StckOprc string `json:"stck_oprc"` // 오늘 시가
			StckHgpr string `json:"stck_hgpr"` // 전일 고가
			StckLwpr string `json:"stck_lwpr"` // 전일 저가
		} `json:"output"`
	}

	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return 0, fmt.Errorf("error unmarshaling response: %v", err)
	}

	// 데이터 추출 및 계산
	if len(result.Output) < 2 {
		return 0, fmt.Errorf("insufficient data in response output")
	}

	stckOprc, err := strconv.Atoi(result.Output[0].StckOprc)
	if err != nil {
		return 0, fmt.Errorf("error converting stck_oprc: %v", err)
	}

	stckHgpr, err := strconv.Atoi(result.Output[1].StckHgpr)
	if err != nil {
		return 0, fmt.Errorf("error converting stck_hgpr: %v", err)
	}

	stckLwpr, err := strconv.Atoi(result.Output[1].StckLwpr)
	if err != nil {
		return 0, fmt.Errorf("error converting stck_lwpr: %v", err)
	}

	targetPrice := float64(stckOprc) + (float64(stckHgpr)-float64(stckLwpr))*0.5
	return targetPrice, nil
}
