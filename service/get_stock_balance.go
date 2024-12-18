package service

import (
	"JSH-Automated-Stock-Trading/config"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/valyala/fasthttp"
)

// 주식 잔고 조회
// https://apiportal.koreainvestment.com/apiservice/apiservice-domestic-stock-order#L_66c61080-674f-4c91-a0cc-db5e64e9a5e6
func GetStockBalance(accessToken string) map[string]string {
	// API 경로 및 URL 설정
	PATH := "uapi/domestic-stock/v1/trading/inquire-balance"
	url := fmt.Sprintf("%s/%s", config.SetConfig.UrlBase, PATH)

	// 요청 헤더 설정
	headers := map[string]string{
		"Content-Type":  "application/json",
		"authorization": "Bearer " + accessToken,
		"appKey":        config.SetConfig.AppKey,
		"appSecret":     config.SetConfig.AppSecret,
		"tr_id":         "VTTC8434R",
		"custtype":      "P",
	}

	// 요청 파라미터 설정
	params := map[string]string{
		"CANO":                  config.SetConfig.Cano,
		"ACNT_PRDT_CD":          config.SetConfig.AcntPrdtCd,
		"AFHR_FLPR_YN":          "N",
		"OFL_YN":                "",
		"INQR_DVSN":             "02",
		"UNPR_DVSN":             "01",
		"FUND_STTL_ICLD_YN":     "N",
		"FNCG_AMT_AUTO_RDPT_YN": "N",
		"PRCS_DVSN":             "01",
		"CTX_AREA_FK100":        "",
		"CTX_AREA_NK100":        "",
	}

	// HTTP 요청 생성
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(url)
	req.Header.SetMethod("GET")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// 파라미터 설정
	args := req.URI().QueryArgs()
	for k, v := range params {
		args.Set(k, v)
	}

	// HTTP 응답 객체 생성
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	// 요청 전송
	if err := fasthttp.Do(req, resp); err != nil {
		log.Fatalf("Error sending request: %v", err)
	}

	// 응답 파싱
	var response struct {
		Output1 []struct {
			ProductName string `json:"prdt_name"`
			ProductCode string `json:"pdno"`
			HoldingQty  string `json:"hldg_qty"`
		} `json:"output1"`
		Output2 []struct {
			EvaluationAmount string `json:"scts_evlu_amt"`
			ProfitLossTotal  string `json:"evlu_pfls_smtl_amt"`
			TotalBalance     string `json:"tot_evlu_amt"`
		} `json:"output2"`
	}

	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		log.Fatalf("Error unmarshaling response: %v", err)
	}

	// 주식 정보 출력 및 맵에 저장
	stockDict := make(map[string]string)
	fmt.Println("==== 주식 보유잔고 ====")
	for _, stock := range response.Output1 {
		if stock.HoldingQty != "0" {
			stockDict[stock.ProductCode] = stock.HoldingQty
			fmt.Printf("%s(%s): %s주\n", stock.ProductName, stock.ProductCode, stock.HoldingQty)
			time.Sleep(100 * time.Millisecond)
		}
	}

	// 평가 금액 정보 출력
	if len(response.Output2) > 0 {
		evaluation := response.Output2[0]
		fmt.Printf("주식 평가 금액: %s원\n", evaluation.EvaluationAmount)
		time.Sleep(100 * time.Millisecond)
		fmt.Printf("평가 손익 합계: %s원\n", evaluation.ProfitLossTotal)
		time.Sleep(100 * time.Millisecond)
		fmt.Printf("총 평가 금액: %s원\n", evaluation.TotalBalance)
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println("=================")

	return stockDict
}
