package trello

import (
	"testing"
)

func Test(t *testing.T) {
	var c = Card{Name: "My Card", Desc: `a long comment and a Startdate: "12.12.2012" `}
	if c.Name != "My Card" {
		t.Fail()
	}
	if c.Startdate() != "12.12.2012" {
		t.Fail()
	}
}
