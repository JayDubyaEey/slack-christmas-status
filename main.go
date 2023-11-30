package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Profile struct {
	StatusText  string `json:"status_text"`
	StatusEmoji string `json:"status_emoji"`
}

type Payload struct {
	Profile Profile `json:"profile"`
}

type SlackClient struct {
	httpClient *http.Client
	token      string
}

func NewSlackClient(token string) *SlackClient {
	return &SlackClient{
		httpClient: &http.Client{},
		token:      token,
	}
}

func (sc *SlackClient) updateStatus(profile Profile) error {

	url := "https://slack.com/api/users.profile.set"
	payload := Payload{Profile: profile}
	jsonPayload, err := json.Marshal(payload)

	if err != nil {

		return fmt.Errorf("error marshaling payload: %w", err)

	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))

	if err != nil {

		return fmt.Errorf("error creating request: %w", err)

	}

	req.Header.Set("Authorization", "Bearer "+sc.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := sc.httpClient.Do(req)

	if err != nil {

		return fmt.Errorf("error making request: %w", err)

	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		return fmt.Errorf("failed to update status. HTTP Status Code: %d", resp.StatusCode)

	}

	return nil
}

func daysUntilChristmas() int {

	currentTime := time.Now()
	christmas := time.Date(currentTime.Year(), time.December, 26, 0, 0, 0, 0, currentTime.Location())
	duration := christmas.Sub(currentTime)
	return int(duration.Hours() / 24)

}

func main() {

	token := os.Getenv("SLACK_AUTH_TOKEN")

	if len(token) == 0 {

		fmt.Println("error: SLACK_AUTH_TOKEN environment variable is not set.")
		return

	}

	sc := NewSlackClient(token)
	days := daysUntilChristmas()
	statusText := fmt.Sprintf("%d days until Christmas", days)
	profile := Profile{StatusText: statusText, StatusEmoji: ":christmas_tree:"}

	if err := sc.updateStatus(profile); err != nil {

		fmt.Println("failed to update Slack status:", err)
		return

	}

	fmt.Println("status updated successfully.")

}
