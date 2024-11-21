package models

type Message struct {
	ID        string `json:"id"`
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	GroupID   string `json:"group_id,omitempty"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}
