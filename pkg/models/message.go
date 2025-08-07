package models

import "time"

type Message struct {
	Id      uint32
	Room    string
	Message string
	User    string
	Sent    time.Time
}
