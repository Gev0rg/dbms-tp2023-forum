package models

import "time"

type Thread struct {
	Id      int       `json:"thread_id" db:"thread_id"`
	Title   string    `json:"title" db:"title"`
	Author  string    `json:"author" db:"author"`
	Forum   string    `json:"forum" db:"forum"`
	Message string    `json:"message" db:"message"`
	Votes   int       `json:"votes" db:"votes"`
	Slug    string    `json:"slug" db:"slug"`
	Created time.Time `json:"created" db:"created"`
}

type CreateThread struct {
	Title   string    `json:"title"`
	Author  string    `json:"author"`
	Message string    `json:"message"`
	Created time.Time `json:"created"`
}

type UpdateThread struct {
	Title   string    `json:"title"`
	Message string    `json:"message"`
}

type GetThreadPostsById struct {
	Id    int64 `json:"id"`
	Limit int64  `json:"limit"`
	Since string `json:"since"`
	Sort  string `json:"sort"`
	Desc  bool   `json:"desc"`
}

type GetThreadPostsBySlug struct {
	Slug  string `json:"slug"`
	Limit int64  `json:"limit"`
	Since string `json:"since"`
	Sort  string `json:"sort"`
	Desc  bool   `json:"desc"`
}
