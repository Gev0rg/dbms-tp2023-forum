package models

import "time"

type Post struct {
	Id       int       `json:"post_id" db:"post_id"`
	Parent   int       `json:"parent" db:"parent"`
	Author   string       `json:"author" db:"author"`
	Message  string    `json:"message" db:"message"`
	IsEdited bool      `json:"is_edited" db:"is_edited"`
	Forum    string    `json:"forum" db:"forum"`
	Thread   int       `json:"thread_id" db:"thread_id"`
	Created  time.Time `json:"created" db:"created"`
}

type CreatePost struct {
	Parent  int    `json:"parent"`
	Author  string `json:"author"`
	Message string `json:"message"`
}

type UpdatePost struct {
	Message  string    `json:"message"`
}

type FullPost struct {
	Post Post `json:"post"`
	Author User `json:"author"`
	Thread Thread `json:"thread"`
	Forum Forum `json:"forum"`
}
