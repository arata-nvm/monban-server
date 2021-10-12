package env

import "os"

func Port() string {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		return "8080"
	}

	return port
}

func GoogleApiToken() string {
	return os.Getenv("GOOGLE_API_TOKEN")
}

func GoogleApiCred() string {
	return os.Getenv("GOOGLE_API_CRED")
}

func EntryLogSID() string {
	return os.Getenv("ENTRY_LOG_SID")
}

func StudentsSID() string {
	return os.Getenv("STUDENTS_SID")
}

func SlackWebhook() string {
	return os.Getenv("SLACK_WEBHOOK")
}
