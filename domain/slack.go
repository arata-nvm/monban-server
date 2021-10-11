package domain

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/arata-nvm/monban/env"
)

func PostToSlack(text string) error {
	webhookUrl := env.SlackWebhook()

	content := map[string]string{"text": text}
	payload, err := json.Marshal(content)
	if err != nil {
		return err
	}

	http.Post(webhookUrl, "application/json", bytes.NewBuffer(payload))

	return nil
}
