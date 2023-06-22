package threads

import "forum/internal/models"

type ThreadUsecase interface {
	CreateThread(thread *models.Thread) (*models.Thread, error)
	GetThreadsByForum(tv *models.ThreadsVars) ([]*models.Thread, error)
	GetUsersByForum(tv *models.ThreadsVars) ([]*models.Thread, error)
	GetThread(slug string, id int64) (*models.Thread, error)
	UpdateThread(thredUpdate *models.ThreadUpdate) (*models.Thread, error)
}
