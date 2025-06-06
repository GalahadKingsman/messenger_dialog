package models

import "time"

type Message struct {
	ID         int
	UserID     int
	DialogID   int
	Text       string
	CreateDate time.Time
}

//Create table messages (
//	id SERIAL PRIMARY KEY,
//	user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
//	dialog_id INTEGER REFERENCES dialogs(id) ON DELETE CASCADE,
//	text TEXT NOT NULL,
//	create_date TIMESTAMP DEFAULT NOW()
//	)
