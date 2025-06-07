package models

import "time"

type Dialog struct {
	ID         int32
	Name       string
	CreateDate time.Time
}

//Create table dialogs (
//	id SERIAL PRIMARY KEY,
//	name TEXT NOT NULL,
//	create_date TIMESTAMP DEFAULT NOW()
//	)
