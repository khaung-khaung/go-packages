package entities

import "time"

type CustomPublish struct {
	Name      string                 `json:"name"`
	Email     string                 `json:"email"`
	Timestamp string                 `json:"timestamp"`
	UserID    string                 `json:"user_id"`
	IsActive  bool                   `json:"is_active"`
	CreatedAt time.Time              `json:"created_at"`
	Metadata  map[string]interface{} `json:"metadata"`
}
