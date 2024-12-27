package service

import (
	"JSH-Automated-Stock-Trading/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Buy(code string, qty int, accessToken string) bool {
	path := "uapi/domestic-stock/v1/trading/order-cash"
	url := fmt.Sprintf("%s/%s", config.SetConfig.UrlBase, path)

	data := map[string]interface{}{
		"CANO":         "your_account_number",     // 계좌번호
		"ACNT_PRDT_CD": "your_account_product_cd", // 계좌상품코드
		"PDNO":         code,                      // 종목 코드
		"ORD_DVSN":     "01",                      // 주문 구분 (시장가)
		"ORD_QTY":      qty,                       // 주문 수량
		"ORD_UNPR":     "0",                       // 주문 단가 (시장가)
	}

	// Serialize request data to JSON
	requestBody, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Error marshalling request data: %v\n", err)
		return false
	}

	// Prepare HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Printf("Error creating HTTP request: %v\n", err)
		return false
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authorization", "Bearer "+accessToken)
	req.Header.Set("appKey", "YourAppKey")
	req.Header.Set("appSecret", "YourAppSecret")
	req.Header.Set("tr_id", "TTTC0802U")
	req.Header.Set("custtype", "P")

	// Send HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending HTTP request: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	// Parse response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return false
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Printf("Error unmarshalling response body: %v\n", err)
		return false
	}

	// Check result code
	if result["rt_cd"] == "0" {
		fmt.Printf("[매수 성공] %v\n", result)
		return true
	} else {
		fmt.Printf("[매수 실패] %v\n", result)
		return false
	}
}
