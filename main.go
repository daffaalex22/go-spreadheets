package main

import (
	"context"
	"fmt"

	freedb "github.com/FreeLeh/GoFreeDB"
	"github.com/FreeLeh/GoFreeDB/google/auth"
)

type Person struct {
	Name       string `db:"Name"`
	Gender     string `db:"Gender"`
	ClassLevel string `db:"Class Level"`
	HomeState  string `db:"Home State"`
	Major      string `db:"Major"`
	ExCul      string `db:"Extracurricular Activity"`
}

func main() {
	// If using Google OAuth2 Flow.
	auth, err := auth.NewOAuth2FromFile(
		"client_secret.json",
		"token.json",
		freedb.FreeDBGoogleAuthScopes,
		auth.OAuth2Config{},
	)
	if err != nil {
		fmt.Printf("Cannot obtain auth from client secret: %v", err)
	}

	store := freedb.NewGoogleSheetRowStore(
		auth,
		"10QcMndAD8j4ms1H2yit-zdgN_a5YBYAt98UH5gIgpWg",
		"Class Data",
		freedb.GoogleSheetRowStoreConfig{
			Columns: []string{
				"Name",
				"Gender",
				"Class Level",
				"Home State",
				"Major",
				"Extracurricular Activity",
			},
		},
	)

	defer store.Close(context.Background())

	// Output variable
	var output []Person

	// Select all columns for all rows
	// err = store.
	// 	Select(&output).
	// 	Exec(context.Background())
	// if err != nil {
	// 	fmt.Printf("Error querying rows: %v", err)
	// }

	// for _, person := range output {
	// 	fmt.Println("Person: ", person.Name)
	// }

	// Filter row
	// err = store.
	// 	Select(&output).
	// 	Where("_rid = ?", 1).
	// 	Exec(context.Background())
	// if err != nil {
	// 	fmt.Printf("Error querying rows: %v", err)
	// }

	// for _, person := range output {
	// 	fmt.Println("Person: ", person.Name)
	// 	fmt.Println("Person: ", person.Gender)
	// 	fmt.Println("Person: ", person.Major)
	// 	fmt.Println("Person: ", person.ClassLevel)
	// 	fmt.Println("Person: ", person.HomeState)
	// 	fmt.Println("Person: ", person.ExCul)
	// }

	// Select rows with sorting/order by
	// ordering := []freedb.ColumnOrderBy{
	// 	{Column: "Name", OrderBy: freedb.OrderByAsc},
	// }

	// err = store.
	// 	Select(&output).
	// 	OrderBy(ordering).
	// 	Exec(context.Background())
	// if err != nil {
	// 	fmt.Printf("Error querying rows: %v", err)
	// }

	// for _, person := range output {
	// 	fmt.Println("Person: ", person.Name)
	// }

	// Select rows with offset and limit
	err = store.
		Select(&output).
		// Offset(2).
		Limit(5).
		Exec(context.Background())
	if err != nil {
		fmt.Printf("Error querying rows: %v", err)
	}

	for _, person := range output {
		fmt.Println("Person: ", person.Name)
	}
}
