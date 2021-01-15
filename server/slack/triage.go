package slack

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"time"

	"github.com/tylerconlee/Deskmate/server/datastore"
)

// Triage outlines the various users that are currently in the
// triage role within Slack. Because of multi-channel triage,
// the channel and user structs are included here to represent the
// multiple channels that could have an active triager
type Triage struct {
	Channel Channel
	User    User
	Started time.Time
}

// Channel struct defines the details about a specific channel in Slack
// that a triage role can be assumed from. The channel ID is Slack's internal
// ID that allows for the channel to be referenced programmatically. The name
// included here allows for the human-friendly channel name to be displayed
// in the frontend UI
type Channel struct {
	Name string
	ID   string
}

// User struct defines the details about a user within Slack that is currently
// in a triage role. The name is the human readable value that can be used as a
// display in the UI/Slack, and the ID is the Slack internal ID that allows for
// that user to be accessed/called programmatically
type User struct {
	Name string
	ID   string
}

// T represents the users that are currently in the triage role
var T []Triage

// GetTriage gets the triage user details for a specific channel and returns it
// in a JSON format so it can be parsed in Slack and the frontend.
// Endpoint: GET /api/triage/{channel-id}
func GetTriage(w http.ResponseWriter, r *http.Request) (n Triage) {
	return
}

func DeleteTriage(w http.ResponseWriter, r *http.Request) (n Triage) {
	u := path.Base(r.RequestURI)
	fmt.Println("Removing active triager for channel: ", u)
	removeTriage(u)
	return
}

// GetAllTriage returns the current triage object to be used on the frontend UI
// to show all users in every active channel that currently has a user in the
// triage role.
// Endpoint: GET /api/triage
func GetAllTriage(w http.ResponseWriter, r *http.Request) {
	// Add LoadTriage to retrieve data from database
	loadTriage()
	t, err := json.Marshal(T)
	fmt.Println(T)
	if err != nil {
		fmt.Println("Error marshalling JSON for config")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(t)

}

func setTriage(channel string, user string) {
	u := getUserInfo(user)
	c := getChannelInfo(channel)

	new := Triage{
		Channel: c,
		User:    u,
		Started: time.Now(),
	}
	saveTriage(new)

	T = append(T, new)
	fmt.Println("Added new triager", T)

}

func saveTriage(t Triage) {
	fmt.Println("Saving triage role to database")
	datastore.SaveTriage(t.User.ID, t.Channel.ID)
}

func loadTriage() {
	fmt.Println("Loading existing triage from database")
	rows := datastore.LoadTriage()
	fmt.Println(rows)

}

func removeTriage(channel string) {
	for i := range T {
		if T[i].Channel.Name == channel {
			T = append(T[:i], T[i+1:]...)
			break
		}
	}
}
