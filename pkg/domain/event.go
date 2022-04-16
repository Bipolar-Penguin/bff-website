package domain

type Event struct {
	GUID    string `json:"guid"`
	Action  string `json:"action"`
	Amount  int    `json:"amount"`
	EventID string `json:"event_id"`
}
