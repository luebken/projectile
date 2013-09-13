package main

import (
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type Member struct {
	Name string `json:"fullName"`
}

type Card struct {
	Name      string   `json:"name"`
	Due       string   `json:"due"`
	Desc      string   `json:"desc"`
	IdMembers []string `json:"idMembers"`
	Members   []Member
}

func (c *Card) LoadMembers() {
	c.Members = make([]Member, len(c.IdMembers))
	for i := 0; i < len(c.IdMembers); i++ {
		body := callTrello("members/" + c.IdMembers[i])
		err := json.Unmarshal(body, &c.Members[i])
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
	}
}

func (c *Card) MembersAsString() string {
	c.LoadMembers()
	if len(c.Members) > 0 {
		result := ""
		for _, member := range c.Members {
			result += member.Name + "; "
		}
		return strings.TrimSuffix(result, "; ")
	} else {
		return "<not assigned>"
	}
}

func (c *Card) Startdate() string {
	re := regexp.MustCompile("Startdate: \"(.*)\"")
	matches := re.FindStringSubmatch(c.Desc)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func callTrello(call string) []byte {
	url := "https://api.trello.com/1/" + call
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
	return body
}

func main() {
	body := callTrello("board/522730ae9504e7ed3d0038e2/cards")

	// JSON parsing
	var cards []Card
	err := json.Unmarshal(body, &cards)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	//fmt.Printf("%+v", cards)

	// Output
	writer := csv.NewWriter(os.Stdout)

	//		  fo, err := os.Create("output.csv")
	//		writer := csv.NewWriter(fo)

	writer.Write([]string{"Arbeitspaket", "Startdatum", "Enddatum", "Kollegen"})
	for _, card := range cards {
		writer.Write([]string{card.Name, card.Due, card.Startdate(), card.MembersAsString()})
	}

	writer.Flush()
}
