package tasks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func IsUserSubscribedToChannel(telegramUsername int) (bool, error) {
	// Cоставляем запрос
	var request struct {
		Telegram_user_id int `json:"telegram_user_id"` 
	}
	request.Telegram_user_id = telegramUsername

	data, err := json.Marshal(request)
	if err != nil {
		return false, err
	}
	// Запрос отправляется на мой личный сервер
	// там развернут hft алгоритм с телеграмм ботом
	// https://t.me/hft_alerts_station
	httpreq, err := http.NewRequest("POST", "http://38.180.148.174:8085/check_subscriber", bytes.NewReader(data))
	if err != nil {
		return false, err
	}
	// Выполняем HTTP-запрос
	cli := http.Client{}

	resp, err := cli.Do(httpreq)
	if err != nil {
		return false, fmt.Errorf("failed to call Telegram API: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("user not subsribed to telegram chat")
	}

	return true, nil
}
