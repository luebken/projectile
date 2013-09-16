package main

import (
	"encoding/csv"
	"encoding/json"
	"github.com/luebken/projectile/trello"
	"log"
	"os"
	"regexp"
)

type Args struct {
	Board     string
	OutputCSV bool
}

func extractCommandLineArgs(args []string) Args {
	boardid := "<undefined>"
	outputcsv := false

	for i := 0; i < len(args); i++ {
		//board-id
		re := regexp.MustCompile("boardid:(.*)")
		matches := re.FindStringSubmatch(args[i])
		if len(matches) > 0 {
			boardid = matches[1]
		}
		//csvoutput
		re = regexp.MustCompile("csvoutput:(.*)")
		matches = re.FindStringSubmatch(args[i])
		if len(matches) > 0 {
			outputcsv = (matches[1] == "true")
		}
	}

	if boardid == "<undefined>" {
		log.Fatal("Usage: projectile boardid:<boardid> csvoutput:false")
	}

	return Args{boardid, outputcsv}
}

func main() {
	args := extractCommandLineArgs(os.Args)

	body := trello.CallTrello("board/" + args.Board + "/cards")

	// JSON parsing
	var cards []trello.Card
	err := json.Unmarshal(body, &cards)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Output
	writer := csv.NewWriter(os.Stdout)

	if args.OutputCSV {
		log.Println("writing to output.csv")
		fo, err := os.Create("output.csv")
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		writer = csv.NewWriter(fo)
	}

	writer.Write([]string{"Arbeitsbereich", "Arbeitspaket", "Status", "Startdatum", "Enddatum", "Kollegen"})
	for _, card := range cards {
		card.LoadList()
		writer.Write([]string{card.LabelsAsString(), card.Name, card.List.Name, card.Due, card.Startdate(), card.MembersAsString()})
	}

	writer.Flush()
}
