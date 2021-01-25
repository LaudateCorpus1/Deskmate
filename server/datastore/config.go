package datastore

import (
	"database/sql"
	"fmt"
)

var configID = 0

// LoadConfig pulls the configuration details from the database and returns a
// pointer to a sql.Row
func LoadConfig() (rows *sql.Row) {
	// Load config from database
	row := db.QueryRow("SELECT slack_url,slack_api,slack_signing from configuration where id = 1 ")
	fmt.Println(row)
	return row
}

// SaveConfig takes a map[string]interface containing configuration details
// and saves the key/values to the database
func SaveConfig(data map[string]interface{}) {
	fmt.Println(data)

	err := db.QueryRow(`INSERT INTO configuration
	(
		id, 
		slack_api, 
		slack_url, 
		slack_signing
	) VALUES (
		1, 
		$1, 
		$2,
		$3
	) ON CONFLICT (id) 
	DO UPDATE SET 
	slack_api = $1, 
	slack_url = $2,
	slack_signing = $3 
	RETURNING id `,
		data["slackapi"],
		data["slackurl"],
		data["slacksigning"]).Scan(&configID)
	if err != nil {
		fmt.Println("error saving configuration into database", err.Error())
	}

}
