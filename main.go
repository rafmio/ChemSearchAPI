package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

const apiUrl = "https://api.wordstat.yandex.net/v4/words"

type StatsResponse struct {
	Data struct {
		Words []struct {
			Text      string `json:"text"`
			Frequency int    `json:"frequency"`
		} `json:"words"`
	} `json:"data"`
}

func main() {
	clientId := os.Getenv("YANDEX_WORDSTAT_CLIENT_ID")
	clientSecret := os.Getenv("YANDEX_WORDSTAT_CLIENT_SECRET")
	code := os.Getenv("AUTHORIZATION_CODE") // Авторизационный код, полученный от пользователя

	fmt.Println(clientId)
	fmt.Println(clientSecret)

	// conf := &oauth2.Config{
	// 	ClientID:     clientId,
	// 	ClientSecret: clientSecret,
	// 	Endpoint: oauth2.Endpoint{
	// 		AuthURL:  "https://oauth.yandex.ru/authorize",
	// 		TokenURL: "https://oauth.yandex.ru/token",
	// 	},
	// 	Scopes: []string{"direct"},
	// }
	conf := &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://oauth.yandex.ru/authorize",
			TokenURL: "https://oauth.yandex.ru/token",
		},
		Scopes: []string{"direct"},
		// Явно укажите передачу параметров в теле запроса
		// RedirectURL: "https://oauth.yandex.ru/verification_code", // Ваш redirect_uri
		RedirectURL: "https://api.wordstat.yandex.net", // Ваш redirect_uri
	}

	token, err := conf.Exchange(context.Background(), code)
	if err != nil {
		log.Fatalf("Failed to exchange auth code for token: %v", err)
	}

	fmt.Println("checkpoint 1")

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, http.DefaultClient)
	httpClient := conf.Client(ctx, token)

	reqBody := map[string]interface{}{
		"queries": []map[string]string{
			{"phrase": "машины"},
		},
	}

	bodyBytes, _ := json.Marshal(reqBody)
	resp, err := httpClient.Post(apiUrl, "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	responseData, _ := io.ReadAll(resp.Body)
	var statsResp StatsResponse
	err = json.Unmarshal(responseData, &statsResp)
	if err != nil {
		log.Fatalf("Error unmarshaling response data: %v", err)
	}

	for _, word := range statsResp.Data.Words {
		fmt.Printf("%s : %d\n", word.Text, word.Frequency)
	}
}
