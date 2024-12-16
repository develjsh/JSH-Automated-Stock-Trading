package service

import (
	"JSH-Automated-Stock-Trading/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

func GetAccessToken() string {
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	body := map[string]string{
		"grant_type": "client_credentials",
		"appkey":     config.SetConfig.AppKey,
		"appsecret":  config.SetConfig.AppSecret,
	}

	bodyJSON, err := json.Marshal(body)
	if err != nil {
		log.Fatalf("Error marshaling request body: %v", err)
	}

	path := "oauth2/tokenP"
	url := fmt.Sprintf("%s/%s", config.SetConfig.UrlBase, path)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	var tokenResponse TokenResponse
	err = json.Unmarshal(respBody, &tokenResponse)
	if err != nil {
		log.Fatalf("Error unmarshaling response body: %v", err)
	}

	return tokenResponse.AccessToken
}
