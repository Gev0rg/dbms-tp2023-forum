package models

type Forum struct {
	Title   string `json:"title"   db:"title"`
	User    string `json:"user"    db:"user"`
	Slug    string `json:"slug"    db:"slug"`
	Posts   int    `json:"posts"   db:"posts"`
	Threads int    `json:"threads" db:"threads"`
}

type CreateForum struct {
	Title string `json:"title"`
	User  string `json:"user"`
	Slug  string `json:"slug"`
}

type GetForumUsers struct {
	Slug  string `json:"slug"`
	Limit int64  `json:"limit"`
	Since string `json:"since"`
	Desc  bool   `json:"desc"`
}

type GetForumThreads struct {
	Slug  string `json:"slug"`
	Limit int64  `json:"limit"`
	Since string `json:"since"`
	Desc  bool   `json:"desc"`
}
