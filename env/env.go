package env

import "os"

// サーバーが使用するポート
func Port() string {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		return "8080"
	}

	return port
}

// Spreadsheetのアクセストークン
func GoogleApiToken() string {
	return os.Getenv("GOOGLE_API_TOKEN")
}

// OAuth 2.0 クライアントの認証情報
func GoogleApiCred() string {
	return os.Getenv("GOOGLE_API_CRED")
}

// 入退室データを記録するスプレッドシートのID
func EntryLogSID() string {
	return os.Getenv("ENTRY_LOG_SID")
}

// 学生データを記録するスプレッドシートのID
func StudentsSID() string {
	return os.Getenv("STUDENTS_SID")
}

// SlackのWebhook URL
func SlackWebhook() string {
	return os.Getenv("SLACK_WEBHOOK")
}
