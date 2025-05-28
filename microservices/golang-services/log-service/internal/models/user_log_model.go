package models

import "time"

type User struct {
	ID    string  `json:"id"`
	Email *string `json:"email,omitempty"`
}

type Log struct {
	Timestamp time.Time              `json:"timestamp"`
	Event     string                 `json:"event"`
	User      *User                  `json:"user"`
	IPAddress string                 `json:"ip_address"`
	SessionID *string                `json:"session_id"`
	Details   map[string]interface{} `json:"details"`
}
