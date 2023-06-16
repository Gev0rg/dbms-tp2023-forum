package usecase

import (
	"context"
	"errors"
	"dbms/internal/models"
	myErrors "dbms/internal/models/errors"
	thread "dbms/internal/thread/repository"
	user "dbms/internal/user/repository"
	forum "dbms/internal/forum/repository"
)

type Usecase interface {
	CreateForum(ctx context.Context, forum models.Forum) (models.Forum, error)
	CreateForumThread(ctx context.Context, slug string, thread models.Thread) (models.Thread, error)

	GetForumBySlug(ctx context.Context, slug string) (models.Forum, error)
	GetForumUsersBySlug(ctx context.Context, slug string, limit int64, since string, desc bool) ([]models.User, error)
	GetForumThreadsBySlug(ctx context.Context, slug string, limit int64, since string, desc bool) ([]models.Thread, error)
}

type usecase struct {
	forumRepository forum.Repository
	userRepository user.Repository
	threadRepository thread.Repository
}

func (u *usecase) Forum(ctx context.Context, forum models.Forum) (models.Forum, error) {
	err := u.userRepository.CheckExistUserByNickname(ctx, forum.User)
	if err != nil {
		return models.Forum{}, err
	}

	existForum, err := u.forumRepository.GetForumBySlug(ctx, forum.Slug)
	if err == nil {
		return existForum, myErrors.ErrForumIsAlreadyExisted
	}

	if err != nil {
		if !errors.Is(err, myErrors.ErrForumNotFound) {
			return models.Forum{}, err
		}
	}

	err = u.forumRepository.CreateForum(ctx, forum)
	if err != nil {
		return models.Forum{}, err
	}

	return forum, nil
}

// func (u *usecase) CreateForumThread(ctx context.Context, slug string, thread models.Thread) (models.Thread, error) {

// }

// func (u *usecase) GetForumBySlug(ctx context.Context, slug string) (models.Forum, error) {
	
// }

// func (u *usecase) GetForumUsersBySlug(ctx context.Context, slug string, limit int64, since string, desc bool) ([]models.User, error) {

// }

// func (u *usecase) GetForumThreadsBySlug(ctx context.Context, slug string, limit int64, since string, desc bool) ([]models.Thread, error) {

// }

// func NewUsecase (forumRepository forum.Repository) Usecase {
// 	return &usecase{forumRepository: forumRepository}
// }
