package domain

import "time"

type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

type Message struct {
	ID             string
	ConversationID string
	Role           Role
	Content        string
	CreatedAt      time.Time
}
