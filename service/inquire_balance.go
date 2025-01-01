package service

import (
	"JSH-Automated-Stock-Trading/config"
	"fmt"
)

// 주식 잔고 조회
func InquireBalance(accessToken string) {
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
}
