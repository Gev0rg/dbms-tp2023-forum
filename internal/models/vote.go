package models

type Vote struct {
	Nickname string `json:"nickname" db:"nickname"`
	Voice    int    `json:"voice" db:"voice"`
}
