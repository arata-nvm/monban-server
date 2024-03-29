package domain

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/arata-nvm/monban/env"
)

// Slackにメッセージを送信する
//
//   - text: メッセージ
func PostToSlack(text string) error {
	webhookUrl := env.SlackWebhook()

	content := map[string]string{"text": text}
	payload, err := json.Marshal(content)
	if err != nil {
		return err
	}

	_, err = http.Post(webhookUrl, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	return nil
}
