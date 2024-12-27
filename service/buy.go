package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type BuyRequest struct {
	CANO         string `json:"CANO"`
	ACNT_PRDT_CD string `json:"ACNT_PRDT_CD"`
	PDNO         string `json:"PDNO"`
	ORD_DVSN     string `json:"ORD_DVSN"`
	ORD_QTY      string `json:"ORD_QTY"`
	ORD_UNPR     string `json:"ORD_UNPR"`
}

func Buy(code string, qty int, accessToken string) bool {
	const URL_BASE = "https://api.example.com"
	const PATH = "uapi/domestic-stock/v1/trading/order-cash"
	url := fmt.Sprintf("%s/%s", URL_BASE, PATH)

	data := BuyRequest{
		CANO:         "12345678",
		ACNT_PRDT_CD: "01",
		PDNO:         code,
		ORD_DVSN:     "01",
		ORD_QTY:      fmt.Sprintf("%d", qty),
		ORD_UNPR:     "0",
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
