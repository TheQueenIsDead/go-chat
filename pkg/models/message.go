package models

import "time"

type Message struct {
	Room    string
	Message string
	User    string
	Sent    time.Time
}
