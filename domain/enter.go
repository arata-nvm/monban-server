package domain

import (
	"fmt"
	"strconv"
	"time"

	"github.com/arata-nvm/monban/database"
	"github.com/arata-nvm/monban/env"
)

func Enter(studentID int) error {
	sheetId := env.EntryLogSid()
	writeRange := "A2"

	studentName, err := FindStudentName(studentID)
	if err != nil {
		return err
	}

	values := []interface{}{timestamp(), studentName}
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

func timestamp() string {
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	now := time.Now().In(jst)
	return now.Format("2006/01/02 15:04:05")
}
