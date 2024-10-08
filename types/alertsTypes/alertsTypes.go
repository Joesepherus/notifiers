package alertsTypes

import "time"

type Alert struct {
	ID           int       `json:"id"`
	TriggerValue float64   `json:"triggerValue"`
	AlertType    string    `json:"alertType"`
	Symbol       string    `json:"symbol"`
	UserID       int       `json:"userID"`
	CreatedAt    time.Time `json:"createdAt"`
	CompletedAt  string    `json:"completedAt"`
}
