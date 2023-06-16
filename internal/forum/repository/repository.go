package repository

import (
	"context"
	"dbms/internal/models"
)

type Repository interface {
	CheckExistForumBySlug(ctx context.Context, slug string) error

	CreateForum(ctx context.Context, forum models.Forum) error

	GetForumBySlug(ctx context.Context, slug string) (models.Forum, error)
	GetForumUsers(ctx context.Context, forumUsersInfo models.GetForumUsers) ([]models.User, error)
	GetForumThreads(ctx context.Context, forumThreadsInfo models.GetForumThreads) ([]models.Thread, error)
}
