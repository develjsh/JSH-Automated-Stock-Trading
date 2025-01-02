package service

import (
	"JSH-Automated-Stock-Trading/config"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/valyala/fasthttp"
)

// 응답 JSON 구조체 정의
type Output2 struct {
	NassAmt string `json:"nass_amt"` // 순자산금액
}

type Response struct {
	Output2 []Output2 `json:"output2"`
}

// 주식 잔고 조회
func InquireBalance(accessToken string) int {
	// URL 및 경로 설정
	path := "uapi/domestic-stock/v1/trading/inquire-balance"
	url := fmt.Sprintf("%s/%s", config.SetConfig.UrlBase, path)

	// 요청 헤더 설정
	headers := map[string]string{
		"Content-Type":  "application/json",
		"authorization": "Bearer " + accessToken,
		"appKey":        config.SetConfig.AppKey,
		"appSecret":     config.SetConfig.AppSecret,
		"tr_id":         "VTTC8434R", // 실전 투자: TTTC8434R
		"custtype":      "P",
	}

	// 요청 파라미터 설정
	params := map[string]string{
		"CANO":                  config.SetConfig.Cano,       // 계좌번호 체계의 앞 8자리
		"ACNT_PRDT_CD":          config.SetConfig.AcntPrdtCd, // 계좌번호 체계의 뒤 2자리
		"AFHR_FLPR_YN":          "N",                         // 시간외단일가 여부 (N: 기본값)
		"OFL_YN":                "",                          // 오프라인 여부 (공란이 기본값)
		"INQR_DVSN":             "01",                        // 조회구분 (01: 대출일별), 02 : 종목별
		"UNPR_DVSN":             "01",                        // 단가구분 (01: 기본값)
		"FUND_STTL_ICLD_YN":     "N",                         // 펀드결제분 포함 여부 (N: 포함하지 않음)
		"FNCG_AMT_AUTO_RDPT_YN": "N",                         // 융자금액 자동상환 여부 (N: 기본값)
		"PRCS_DVSN":             "00",                        // 처리구분 (00: 전일매매포함), 01 : 전일매매미포함
		"CTX_AREA_FK100":        "",                          // 연속조회 검색 조건 (최초 조회시 공란)
		"CTX_AREA_NK100":        "",                          // 연속조회 키 (최초 조회시 공란)
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

	// 응답 본문
	body := resp.Body()

	// 응답 JSON 파싱
	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		SendMessage("InquireBalance error parsing JSON", config.SetConfig.DiscordWebhookUrl)
		return -1
	}

	// 응답 상태 코드
	statusCode := resp.StatusCode()
	if statusCode != 200 {
		SendMessage("InquireBalance 응답 상태 !200", config.SetConfig.DiscordWebhookUrl)
		return -1
	}
	nassAmt, err := strconv.Atoi(response.Output2[0].NassAmt)
	if err != nil {
		SendMessage("InquireBalance eerror converting nass_amt to int", config.SetConfig.DiscordWebhookUrl)
		return -1
	}

	return nassAmt
}
