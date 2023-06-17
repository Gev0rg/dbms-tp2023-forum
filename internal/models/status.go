package models

type Status struct {
	Users   int `json:"users"`
	Forums  int `json:"forums"`
	Posts   int `json:"posts"`
	Threads int `json:"threads"`
}
