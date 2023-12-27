package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gopkg.in/Iwark/spreadsheet.v2"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func main() {
	data, err := os.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read Client Secret: %v", err)
	}

	conf, err := google.ConfigFromJSON(data, spreadsheet.Scope)
	if err != nil {
		log.Fatalf("Unable to create Client Secret: %v", err)
	}

	client := getClient(conf)
	service := spreadsheet.NewServiceWithClient(client)

	spreadsheet, err := service.FetchSpreadsheet("10QcMndAD8j4ms1H2yit-zdgN_a5YBYAt98UH5gIgpWg")
	if err != nil {
		log.Fatalf("Unable to fetch Sheets: %v", err)
	}

	sheet, err := spreadsheet.SheetByIndex(0)
	if err != nil {
		log.Fatalf("Unable to get a sheet: %v", err)
	}

	for _, row := range sheet.Rows {
		for _, cell := range row {
			fmt.Println(cell.Value)
		}
	}

	// Update cell content
	sheet.Update(0, 0, "hogehoge")

	// Make sure call Synchronize to reflect the changes
	err = sheet.Synchronize()
	if err != nil {
		log.Fatalf("Unable to update a cell: %v", err)
	}
}
