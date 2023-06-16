package repository

import (
	"context"
	"database/sql"
	"dbms/internal/models"
	myErrors "dbms/internal/models/errors"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	CheckExistForumBySlug(ctx context.Context, slug string) error

	CreateForum(ctx context.Context, forum models.Forum) error

	GetForumBySlug(ctx context.Context, slug string) (models.Forum, error)
	GetForumUsers(ctx context.Context, forumUsersInfo models.GetForumUsers) ([]models.User, error)
	GetForumThreads(ctx context.Context, forumThreadsInfo models.GetForumThreads) ([]models.Thread, error)
}

type repository struct {
	db *sqlx.DB
}

func (r *repository) CheckExistForumBySlug(ctx context.Context, slug string) error {
	var exists bool

	err := r.db.GetContext(ctx, &exists, "SELECT EXISTS(SELECT 1 FROM forums WHERE LOWER(slug) = LOWER($1))", slug)
	if err!= nil {
        return err
    }

	if !exists {
		return myErrors.ErrForumNotFound
	}

	return nil
}

func (r *repository) CreateForum(ctx context.Context, forum models.Forum) error {
	_, err := r.db.NamedExecContext(ctx, `INSERT INTO forums VALUES (:title, :user, :slug, :posts, :threads)`, forum)
	return err
}

func (r *repository) GetForumBySlug(ctx context.Context, slug string) (models.Forum, error) {
	var forum models.Forum

	err := r.db.GetContext(ctx, &forum, `SELECT * FROM forums WHERE LOWER(slug) = LOWER($1)`, slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Forum{}, myErrors.ErrForumNotFound
		}

		return models.Forum{}, err
	}

	return forum, nil
}

func (r *repository) GetForumUsers(ctx context.Context, forumUsersInfo models.GetForumUsers) ([]models.User, error) {
	var users []models.User
	var err error

	if forumUsersInfo.Desc {
		err = r.db.SelectContext(ctx, &users, `SELECT u.nickname, u.fullname, u.about, u.email
			FROM user_forums u WHERE u.forum = $1 AND u.nickname != $2
			ORDER BY LOWER(u.nickname) DESC LIMIT $3`, forumUsersInfo.Slug, forumUsersInfo.Since, forumUsersInfo.Limit)
	} else {
		err = r.db.SelectContext(ctx, &users, `SELECT u.nickname, u.fullname, u.about, u.email
			FROM user_forums u WHERE u.forum = $1 AND u.nickname != $2
			ORDER BY LOWER(u.nickname) LIMIT $3`, forumUsersInfo.Slug, forumUsersInfo.Since, forumUsersInfo.Limit)
	}

	if err != nil {
		return []models.User{}, err
	}

	return users, nil
}

func (r *repository) GetForumThreads(ctx context.Context, forumThreadsInfo models.GetForumThreads) ([]models.Thread, error) {
	var threads []models.Thread
	var err error

	if forumThreadsInfo.Desc {
		err = r.db.SelectContext(ctx, &threads, `SELECT t.id, t.title, t.author, t.forum, t.message, t.votes, t.slug, t.created
			FROM threads AS t
         	LEFT JOIN forums f ON t.forum = f.slug
			WHERE f.slug = $1 AND t.created >= $2
			ORDER BY t.created DESC LIMIT $3`,
			forumThreadsInfo.Slug, forumThreadsInfo.Since, forumThreadsInfo.Limit)
	} else {
		err = r.db.SelectContext(ctx, &threads, `SELECT t.id, t.title, t.author, t.forum, t.message, t.votes, t.slug, t.created
			FROM threads AS t
         	LEFT JOIN forums f ON t.forum = f.slug
			WHERE f.slug = $1 AND t.created >= $2
			ORDER BY t.created LIMIT $3`,
			forumThreadsInfo.Slug, forumThreadsInfo.Since, forumThreadsInfo.Limit)
	}

	if err != nil {
		return []models.Thread{}, err
	}

	return threads, nil
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}
