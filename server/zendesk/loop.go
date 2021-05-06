package zendesk

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/shurcooL/graphql"
)

var active = false
var activeTickets []Ticket
var lastRan time.Time

func Connect(host string) {
	var a string
	if os.Getenv("APP_ENV") == "development" {
		host = strings.Replace(host, "3", "6", 1)

		a = fmt.Sprintf("%squery", host)
	} else {
		a = fmt.Sprintf("%szendesk/query", host)
	}
	client = graphql.NewClient(a, nil)

	variables = map[string]interface{}{
		"user":   c.user,
		"apikey": c.apikey,
		"url":    c.url,
	}
	if !active {
		RunTimer(time.Minute)
	}
}

func RunTimer(interval time.Duration) {
	t := time.NewTicker(interval)
	active = true
	for {
		iteration(t, interval)
		<-t.C

	}
}

func iteration(t *time.Ticker, interval time.Duration) {

	lastRan = time.Now()
	getAllTickets()
	for _, ticket := range activeTickets {
		notify := checkTag(ticket)
		if notify != nil {

			for _, t := range notify {

				switch t.notificationType {
				case "breaches":
					getUser(&ticket)
					sendSLANotification(ticket, t.channel, t.tag)
				case "new":
					sendNewNotification(ticket, t.channel, t.tag)
				case "updates":
					getUser(&ticket)
					sendUpdatedNotification(ticket, t.channel, t.tag)
				}
			}
		}

	}

	activeTickets = nil
	<-t.C
}
