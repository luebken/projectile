package main

import (
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

type Card struct {
	Name string `json:"name"`
	Due  string `json:"due"`
	Desc string `json:"desc"`
}

func (c *Card) startDate() string {
	re := regexp.MustCompile("Startdate: \"(.*)\"")
	matches := re.FindStringSubmatch(c.Desc)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func main() {
	// Get
	url := "https://api.trello.com/1/board/522730ae9504e7ed3d0038e2/cards"
	key := os.Getenv("TRELLO_API_KEY")
	token := os.Getenv("TRELLO_API_TOKEN")
	url = url + "?key=" + key + "&token=" + token
	log.Println("Getting: " + url)
	res, err := http.Get(url)

	// Error handling
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 300 {
		log.Fatalf("Status: %d, Body: %s", res.StatusCode, body)
	}
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// JSON parsing
	var cards []Card
	err = json.Unmarshal(body, &cards)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	//fmt.Printf("%+v", cards)

	// Output
	writer := csv.NewWriter(os.Stdout)

	//		  fo, err := os.Create("output.csv")
	//		writer := csv.NewWriter(fo)

	writer.Write([]string{"card.Name", "card.Due", "card.Startdate"})
	for _, card := range cards {
		writer.Write([]string{card.Name, card.Due, card.startDate()})
	}

	writer.Flush()
}
