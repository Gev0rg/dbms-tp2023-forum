package threads

import "forum/internal/models"

type ThreadRepository interface {
	InsertThread(thread *models.Thread) error
	SelectThreadBySlug(slug string) (*models.Thread, error)
	SelectThreadsByForum(tv *models.ThreadsVars) ([]*models.Thread, error)
	SelectUsersByForum(tv *models.ThreadsVars) ([]*models.Thread, error)
	SelectThread(slug string, id int64) (*models.Thread, error)
	UpdateThread(threadUpdate *models.ThreadUpdate) (*models.Thread, error)
}
