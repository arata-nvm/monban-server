package database

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/arata-nvm/monban/env"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var service *sheets.Service

func Initialize() error {
	cred := env.GoogleApiCred()
	config, err := google.ConfigFromJSON([]byte(cred), "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return err
	}

	token := env.GoogleApiToken()
	tok := &oauth2.Token{}
	err = json.NewDecoder(strings.NewReader(token)).Decode(tok)
	if err != nil {
		return err
	}

	ctx := context.Background()
	client := config.Client(context.Background(), tok)
	service, err = sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}

	return nil
}

// スプレッドシートから値を取得する
//
//   - sheetID: スプレッドシートのID
//   - readRange: 値を取得する範囲
func GetValues(sheetID string, readRange string) ([][]interface{}, error) {
	resp, err := service.Spreadsheets.Values.Get(sheetID, readRange).Do()
	if err != nil {
		return nil, err
	}

	return resp.Values, nil
}

// スプレッドシートに行を追加する
//
//   - sheetID: スプレッドシートのID
//   - values: 追加する行のデータ
func AppendValues(sheetID string, writeRange string, values []interface{}) error {
	rb := &sheets.ValueRange{
		Values: [][]interface{}{values},
	}

	_, err := service.Spreadsheets.Values.Append(
		sheetID, writeRange, rb,
	).ValueInputOption("USER_ENTERED").
		InsertDataOption("INSERT_ROWS").
		Do()
	if err != nil {
		return err
	}

	return nil
}
