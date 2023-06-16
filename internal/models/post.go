package model

import "time"

type Post struct {
	Id       int       `json:"post_id" db:"post_id"`
	Parent   int       `json:"parent" db:"parent"`
	Author   int       `json:"author" db:"author"`
	Message  string    `json:"message" db:"message"`
	IsEdited bool      `json:"is_edited" db:"is_edited"`
	Forum    string    `json:"forum" db:"forum"`
	Thread   int       `json:"thread_id" db:"thread_id"`
	Created  time.Time `json:"created" db:"created"`
}
