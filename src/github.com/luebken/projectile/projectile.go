package main

import (
	"encoding/csv"
	"encoding/json"
	"github.com/luebken/projectile/trello"
	"log"
	"os"
)

//TODO trello board as commandline parameter
func main() {

	body := trello.CallTrello("board/522730ae9504e7ed3d0038e2/cards")

	// JSON parsing
	var cards []trello.Card
	err := json.Unmarshal(body, &cards)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	//fmt.Printf("%+v", cards)

	// Output
	writer := csv.NewWriter(os.Stdout)

	//		  fo, err := os.Create("output.csv")
	//		writer := csv.NewWriter(fo)

	writer.Write([]string{"Arbeitsbereich", "Arbeitspaket", "Startdatum", "Enddatum", "Kollegen"})
	for _, card := range cards {
		writer.Write([]string{card.LabelsAsString(), card.Name, card.Due, card.Startdate(), card.MembersAsString()})
	}

	writer.Flush()
}
