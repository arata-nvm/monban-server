package domain

import (
	"fmt"
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

	SendNotify(now, studentID, studentName)

	sheetId := env.EntryLogSid()
	writeRange := "A2"
	values := []interface{}{now, studentID, studentName}
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

func SendNotify(now string, studentId int, studentName string) error {
	studentIds, err := FindTodaysStudents()
	if err != nil {
		return err
	}

	if len(studentIds) == 0 {
		PostToSlack("🔓")
	}

	numOfRecords := 0
	for _, id := range studentIds {
		if id == studentId {
			numOfRecords += 1
		}
	}

	if numOfRecords%2 == 0 {
		PostToSlack(fmt.Sprintf("%s\n%s さんがログインしました。", now, studentName))
	} else {
		PostToSlack(fmt.Sprintf("%s\n%s さんがログアウトしました。", now, studentName))
	}

	return nil
}

func FindTodaysStudents() ([]int, error) {
	sheetId := env.EntryLogSid()
	readRange := "A2:B"
	values, err := database.GetValues(sheetId, readRange)
	if err != nil {
		return nil, err
	}

	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	now := time.Now().In(jst)

	var studentIds []int
	for i := range values {
		row := values[len(values)-i-1]

		enteredAt, err := time.Parse(TIMESTAMP_FORMAT, row[0].(string))
		if err != nil {
			return nil, err
		}

		if enteredAt.Year() != now.Year() || enteredAt.Month() != now.Month() || enteredAt.Day() != now.Day() {
			break
		}

		studentId, err := strconv.Atoi(row[1].(string))
		if err != nil {
			return nil, err
		}
		studentIds = append(studentIds, studentId)
	}

	return studentIds, nil
}

const TIMESTAMP_FORMAT string = "2006/01/02 15:04:05"

func timestamp() string {
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	now := time.Now().In(jst)
	return now.Format(TIMESTAMP_FORMAT)
}
