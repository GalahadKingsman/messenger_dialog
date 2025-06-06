package models

import "time"

type DialogInfo struct {
	ID           int32
	PeerID       int32
	PeerLogin    string
	LastMessage  string
	LastActivity time.Time
}
