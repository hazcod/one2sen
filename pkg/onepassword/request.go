package onepassword

type eventRequest struct {
	Limit     uint32 `json:"limit,omitempty"`
	Cursor    string `json:"cursor,omitempty"`
	StartTime string `json:"start_time,omitempty"`
	EndTime   string `json:"end_time,omitempty"`
}
