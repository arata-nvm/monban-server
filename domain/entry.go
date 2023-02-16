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
	// 学生がログインした
	EVENT_ENTRY EventType = iota

	// 学生が最初にログインした（🔓の通知用）
	EVENT_FIRST_ENTRY

	// 学生がログアウトした
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

// 学籍番号に対応する学生の名前を取得する
//
//   - studentID: 学籍番号
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

// 現在ログインしている学生の学籍番号を取得する
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

// 2つのtime.Time型のデータについて、日付が同じかを返す
//
//   - t1: 比較する日時
//   - t2: 比較する日時
func DateEquals(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
}

// 同じ学生が直前（10秒前まで）にログインしたかを返す。
// カードリーダーの特性上、短時間で複数回ログインのリクエストが飛ぶことがあるためその確認に用いる。
//
// - studentID: 学生の学籍番号
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

// 学生がログインしたのかログアウトしたのかを返す
//
// - activeStudents: 現在ログインしている学生の学籍番号
// - studentID: 学生の学籍番号
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

// イベントの種別に応じた通知を送信する
//
//   - now: 現在時刻
//   - studentName: 学生の名前
//   - event: イベントの種別
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

// イベントの種別に応じたデータを記録に追加する
//
//   - now: 現在時刻
//   - studentID: 学籍番号
//   - studentName: 学生の名前
//   - typ: イベントの種別
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

// 現在時刻を文字列で返す
func timestamp() string {
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	now := time.Now().In(jst)
	return now.Format(TIMESTAMP_FORMAT)
}

// 配列に特定の値が何個含まれているかを返す
//
//   - arr: 配列
//   - value: カウントする値
func count(arr []int, value int) int {
	cnt := 0
	for _, v := range arr {
		if v == value {
			cnt += 1
		}
	}
	return cnt
}
