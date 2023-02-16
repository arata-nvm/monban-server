package domain

import (
	"fmt"
	"strconv"
	"time"

	"github.com/arata-nvm/monban/database"
	"github.com/arata-nvm/monban/env"
)

type EventType int

const (
	EVENT_ENTRY EventType = iota
	EVENT_FIRST_ENTRY
	EVENT_EXIT
)

func Entry(studentID int) error {
	isDuplicated, err := IsDuplicated(studentID)
	if err != nil {
		return err
	}
	if isDuplicated {
		return nil
	}

	studentName, err := FindStudentName(studentID)
	if err != nil {
		return err
	}

	activeStudents, err := FindActiveStudents()
	if err != nil {
		return err
	}

	now := timestamp()
	event := DetermineEventType(activeStudents, studentID)
	err = PostMessage(now, studentName, event)
	if err != nil {
		return err
	}

	return AppendLog(now, studentID, studentName, event)
}

func FindStudentName(studentID int) (string, error) {
	sheetID := env.StudentsSID()
	readRange := "C2:D"
	values, err := database.GetValues(sheetID, readRange)
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

func FindActiveStudents() ([]int, error) {
	sheetID := env.EntryLogSID()
	readRange := "A2:B"
	values, err := database.GetValues(sheetID, readRange)
	if err != nil {
		return nil, err
	}

	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	now := time.Now().In(jst)

	var studentIDs []int
	for i := range values {
		row := values[len(values)-i-1]

		enteredAt, err := time.Parse(TIMESTAMP_FORMAT, row[0].(string))
		if err != nil {
			return nil, err
		}

		if !DateEquals(enteredAt, now) {
			break
		}

		studentID, err := strconv.Atoi(row[1].(string))
		if err != nil {
			return nil, err
		}
		studentIDs = append(studentIDs, studentID)
	}

	return studentIDs, nil
}

func DateEquals(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
}

func IsDuplicated(studentID int) (bool, error) {
	sheetID := env.EntryLogSID()
	readRange := "A2:B"
	values, err := database.GetValues(sheetID, readRange)
	if err != nil {
		return false, err
	}

	row := values[len(values)-1]
	lastStudentID, err := strconv.Atoi(row[1].(string))
	if err != nil {
		return false, err
	}

	if lastStudentID != studentID {
		return false, nil
	}

	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	enteredAt, err := time.ParseInLocation(TIMESTAMP_FORMAT, row[0].(string), jst)
	if err != nil {
		return false, err
	}

	now := time.Now()
	duration := now.In(jst).Sub(enteredAt)

	if duration.Seconds() < 10 {
		return true, nil
	}

	return false, nil

}

func DetermineEventType(activeStudents []int, studentID int) EventType {
	if len(activeStudents) == 0 {
		return EVENT_FIRST_ENTRY
	}

	numOfRecords := count(activeStudents, studentID)
	if numOfRecords%2 == 0 {
		return EVENT_ENTRY
	} else {
		return EVENT_EXIT
	}
}

func PostMessage(now string, studentName string, event EventType) error {
	switch event {
	case EVENT_FIRST_ENTRY:
		if err := PostToSlack("🔓"); err != nil {
			return err
		}
		fallthrough
	case EVENT_ENTRY:
		return PostToSlack(fmt.Sprintf("%s\n%s さんがログインしました。", now, studentName))
	case EVENT_EXIT:
		return PostToSlack(fmt.Sprintf("%s\n%s さんがログアウトしました。", now, studentName))
	}
	return nil
}

func AppendLog(now string, studentID int, studentName string, typ EventType) error {
	sheetID := env.EntryLogSID()
	writeRange := "A2"

	var typStr string
	switch typ {
	case EVENT_ENTRY, EVENT_FIRST_ENTRY:
		typStr = "入室"
	case EVENT_EXIT:
		typStr = "退室"
	}

	values := []interface{}{now, studentID, studentName, typStr}
	err := database.AppendValues(sheetID, writeRange, values)
	if err != nil {
		return err
	}

	return nil
}

const TIMESTAMP_FORMAT string = "2006/01/02 15:04:05"

func timestamp() string {
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	now := time.Now().In(jst)
	return now.Format(TIMESTAMP_FORMAT)
}

func count(arr []int, value int) int {
	cnt := 0
	for _, v := range arr {
		if v == value {
			cnt += 1
		}
	}
	return cnt
}
