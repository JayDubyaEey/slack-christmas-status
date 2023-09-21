package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// Profile holds the status text and emoji
type Profile struct {
	StatusText  string `json:"status_text"`
	StatusEmoji string `json:"status_emoji"`
}

// Payload wraps the profile data for API request
type Payload struct {
	Profile Profile `json:"profile"`
}

func DaysUntilChristmas() int {
	// Get the current time
	currentTime := time.Now()

	// Create a new time object for December 25th of the current year
	christmas := time.Date(currentTime.Year(), time.December, 25, 0, 0, 0, 0, currentTime.Location())

	// Calculate the duration between now and Christmas
	duration := christmas.Sub(currentTime)

	// Convert duration to days
	days := int(duration.Hours() / 24)

	return days
}

func main() {
	// Replace this with your Bot User OAuth Token
	token := os.Getenv("AUTH_TOKEN")

	client := &http.Client{}
	url := "https://slack.com/api/users.profile.set"

	daysUntilChristmas := DaysUntilChristmas()

	statusTextString := fmt.Sprintf("%d days until Christmas", daysUntilChristmas)

	// Create payload
	profile := Profile{StatusText: statusTextString, StatusEmoji: ":christmas_tree:"}
	
	payload := Payload{Profile: profile}

	// Convert payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshaling payload:", err)
		return
	}

	// Create a new request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	// Make the API call
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	// Check for successful response
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Status updated successfully.", resp.Body)
	} else {
		fmt.Printf("Failed to update status. HTTP Status Code: %d\n", resp.StatusCode)
	}
}
