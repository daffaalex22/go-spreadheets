package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func oldMain() {
	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	// Prints the names and majors of students in a sample spreadsheet:
	// https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit
	spreadsheetId := "10QcMndAD8j4ms1H2yit-zdgN_a5YBYAt98UH5gIgpWg"

	readRange := "Class Data!A2:E"
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		fmt.Println("Name, Major:")
		for _, row := range resp.Values {
			// Print columns A and E, which correspond to indices 0 and 4.
			fmt.Printf("%s, %s\n", row[0], row[4])
		}
	}

	// Update the names two students in a sample spreadsheet
	updateRange := "Class Data!A2:A3"
	var updateValueRange sheets.ValueRange

	updateValues := [][]interface{}{
		{"Alexander Collins"},
		{"Anthony Stephen"},
	}
	updateValueRange.Values = updateValues

	updateResponse, err := srv.Spreadsheets.Values.
		Update(spreadsheetId, updateRange, &updateValueRange).
		ValueInputOption("RAW").
		IncludeValuesInResponse(true).
		Do()
	if err != nil {
		log.Fatalf("Unable to update data from sheet: %v", err)
	}

	if len(updateResponse.UpdatedData.Values) == 0 {
		fmt.Println("No data updated.")
	} else {
		fmt.Println("New Name:")
		for _, val := range updateResponse.UpdatedData.Values {
			fmt.Printf("%s \n", val)
		}
	}

	// Append the names two students in a sample spreadsheet
	appendRange := "Class Data!A2:A3"
	var appendValueRange sheets.ValueRange

	appendValues := [][]interface{}{
		{"Alexander Collins"},
		{"Anthony Stephen"},
	}
	appendValueRange.Values = appendValues

	appendResponse, err := srv.Spreadsheets.Values.
		Append(spreadsheetId, appendRange, &appendValueRange).
		ValueInputOption("RAW").
		IncludeValuesInResponse(true).
		Do()
	if err != nil {
		log.Fatalf("Unable to append data from sheet: %v", err)
	}

	if len(appendResponse.Updates.UpdatedData.Values) == 0 {
		fmt.Println("No data appended.")
	} else {
		fmt.Println("Appended Name:")
		for _, val := range appendResponse.Updates.UpdatedData.Values {
			fmt.Printf("%s \n", val)
		}
	}

	// Delete the previously appended student names in an example spreadsheet
	clearRange := "Class Data!A32:A33"
	clearResponse, err := srv.Spreadsheets.Values.
		Clear(spreadsheetId, clearRange, &sheets.ClearValuesRequest{}).
		Do()
	if err != nil {
		log.Fatalf("Unable to clear data from sheet: %v", err)
	}

	fmt.Printf("%s \n", clearResponse.ClearedRange)
}
