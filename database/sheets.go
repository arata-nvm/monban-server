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
	client := config.Client(context.Background(), tok)

	ctx := context.Background()
	service, err = sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}

	return nil
}

func GetValues(sheetId string, readRange string) ([][]interface{}, error) {
	resp, err := service.Spreadsheets.Values.Get(sheetId, readRange).Do()
	if err != nil {
		return nil, err
	}

	return resp.Values, nil
}

func AppendValues(sheetId string, writeRange string, values []interface{}) error {
	valueInputOption := "USER_ENTERED"
	insertDataOption := "INSERT_ROWS"
	rb := &sheets.ValueRange{
		Values: [][]interface{}{values},
	}

	_, err := service.Spreadsheets.Values.Append(
		sheetId, writeRange, rb,
	).ValueInputOption(valueInputOption).
		InsertDataOption(insertDataOption).
		Do()
	if err != nil {
		return err
	}

	return nil
}
