package domain

import (
	"time"

	"github.com/arata-nvm/monban/database"
	"github.com/arata-nvm/monban/env"
)

func Enter(studentID int) error {
	sheetId := env.EntryLogSid()
	writeRange := "A2"
	values := []interface{}{timestamp(), studentID}
	err := database.AppendValues(sheetId, writeRange, values)
	if err != nil {
		return err
	}

	return nil
}

func timestamp() string {
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	now := time.Now().In(jst)
	return now.Format("2006/01/02 15:04:05")
}
