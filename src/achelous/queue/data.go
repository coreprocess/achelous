package queue

import "time"

type MessageStatus int8

const (
	MessageStatusPreparing MessageStatus = 0
	MessageStatusQueued    MessageStatus = 1
)

type Participant struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type MessageMeta struct {
	Participants struct {
		From *Participant  `json:"from"`
		To   []Participant `json:"to"`
	} `json:"participants"`
	Subject   string    `json:"subject"`
	Timestamp time.Time `json:"timestamp"`
}

type AttachmentMeta struct {
	Id      string `json:"id"`
	Type    string `json:"type"`
	Charset string `json:"charset"`
	Name    string `json:"name"`
}