package trello

import (
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

type Label struct {
	Color string `json:"color"`
	Name  string `json:"name"`
}

type Card struct {
	Name      string   `json:"name"`
	Due       string   `json:"due"`
	Desc      string   `json:"desc"`
	IdMembers []string `json:"idMembers"`
	Labels    []Label  `json:"labels"`
	Members   []Member
}

func (c *Card) LoadMembers() {
	c.Members = make([]Member, len(c.IdMembers))
	for i := 0; i < len(c.IdMembers); i++ {
		body := CallTrello("members/" + c.IdMembers[i])
		err := json.Unmarshal(body, &c.Members[i])
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
	}
}

func (c *Card) LabelsAsString() string {
	result := ""
	for _, label := range c.Labels {
		result += label.Name + "; "
	}
	return strings.TrimSuffix(result, "; ")
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
	re := regexp.MustCompile("Startdate: \"(.*)\"") //TODO only get date
	matches := re.FindStringSubmatch(c.Desc)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

var callCache = make(map[string][]byte)

func CallTrello(call string) []byte {
	result, cacheHit := callCache[call]
	if cacheHit {
		log.Println("Cache hit for " + call)
		return result
	}
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
	callCache[call] = body
	return body
}
