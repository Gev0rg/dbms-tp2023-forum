package models

type User struct {
	Nickname string `json:"nickname" db:"nickname"`
	FullName string `json:"fullname" db:"fullname"`
	About    string `json:"about"    db:"about"`
	Email    string `json:"email"    db:"email"`
}

type UpdateUser struct {
	FullName string `json:"fullname"`
	About    string `json:"about"`
	Email    string `json:"email"`
}
