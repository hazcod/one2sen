package onepassword

const (
	eventsURL = "https://events.1password.com"
)

type eventRequest struct {
	Limit     uint32 `json:"limit"`
	Cursor    string `json:"cursor,omitempty"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}
