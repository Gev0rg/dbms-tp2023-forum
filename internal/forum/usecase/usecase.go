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
	GetForumUsersBySlug(ctx context.Context, forumThreadsInfo models.GetForumUsers) ([]models.User, error)
	GetForumThreadsBySlug(ctx context.Context, forumThreadsInfo models.GetForumThreads) ([]models.Thread, error)
}

type usecase struct {
	forumRepository forum.Repository
	userRepository user.Repository
	threadRepository thread.Repository
}

func (u *usecase) CreateForum(ctx context.Context, forum models.Forum) (models.Forum, error) {
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

func (u *usecase) CreateForumThread(ctx context.Context, slug string, thread models.Thread) (models.Thread, error) {
	err := u.userRepository.CheckExistUserByNickname(ctx, thread.Author)
	if err!= nil {
        return models.Thread{}, err
    }

	err = u.forumRepository.CheckExistForumBySlug(ctx, slug)
	if err!= nil {
        return models.Thread{}, err
    }

	existThread, err := u.threadRepository.GetThreadBySlug(ctx, slug)
	if err == nil {
        return existThread, myErrors.ErrThreadIsAlreadyExisted
    }

	if err != nil {
		if !errors.Is(err, myErrors.ErrThreadNotFound) {
            return models.Thread{}, err
        }
	}
    
	thread, err = u.threadRepository.CreateThread(ctx, thread)
	if err!= nil {
        return models.Thread{}, err
    }

	return thread, nil
}

func (u *usecase) GetForumBySlug(ctx context.Context, slug string) (models.Forum, error) {
	forum, err := u.forumRepository.GetForumBySlug(ctx, slug)
	if err != nil {
		return models.Forum{}, err
	}
		
	return forum, nil
}

func (u *usecase) GetForumUsersBySlug(ctx context.Context, forumUsersInfo models.GetForumUsers) ([]models.User, error) {
	err := u.forumRepository.CheckExistForumBySlug(ctx, forumUsersInfo.Slug)
	if err!= nil {
        return []models.User{}, err
    }

	users, err := u.forumRepository.GetForumUsers(ctx, forumUsersInfo)
	if err!= nil {
        return []models.User{}, err
    }

	return users, nil
}

func (u *usecase) GetForumThreadsBySlug(ctx context.Context, forumThreadsInfo models.GetForumThreads) ([]models.Thread, error) {
	err := u.forumRepository.CheckExistForumBySlug(ctx, forumThreadsInfo.Slug)
	if err!= nil {
        return []models.Thread{}, err
    }

	threads, err := u.forumRepository.GetForumThreads(ctx, forumThreadsInfo)
	if err!= nil {
        return []models.Thread{}, err
    }

	return threads, nil
}

func NewUsecase(forumRepository forum.Repository, userRepository user.Repository, threadRepository thread.Repository) Usecase {
	return &usecase{forumRepository: forumRepository, userRepository: userRepository, threadRepository: threadRepository}
}
