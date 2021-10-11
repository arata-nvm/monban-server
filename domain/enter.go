package domain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/arata-nvm/monban/database"
	"github.com/arata-nvm/monban/env"
)

func Enter(studentID int) error {
	now := timestamp()
	studentName, err := FindStudentName(studentID)
	if err != nil {
		return err
	}

	SendNotify(now, studentName)

	sheetId := env.EntryLogSid()
	writeRange := "A2"
	values := []interface{}{now, studentName}
	err = database.AppendValues(sheetId, writeRange, values)
	if err != nil {
		return err
	}

	return nil
}

func FindStudentName(studentID int) (string, error) {
	sheetId := env.StudentsSid()
	readRange := "C2:D"
	values, err := database.GetValues(sheetId, readRange)
	if err != nil {
		return "", err
	}

	name := fmt.Sprintf("未登録(%d)", studentID)

	for i := range values {
		row := values[len(values)-i-1]
		id, err := strconv.Atoi(row[1].(string))
		if err != nil {
			return "", err
		}

		if id == studentID {
			name = row[0].(string)
			break
		}
	}

	return name, nil
}

func SendNotify(now string, studentName string) error {
	webhookUrl := env.SlackWebhook()

	content := map[string]string{
		"text": fmt.Sprintf("%s\n%s さんがログインしました。", now, studentName),
	}
	payload, err := json.Marshal(content)
	if err != nil {
		return err
	}

	http.Post(webhookUrl, "application/json", bytes.NewBuffer(payload))

	return nil
}

func timestamp() string {
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	now := time.Now().In(jst)
	return now.Format("2006/01/02 15:04:05")
}
