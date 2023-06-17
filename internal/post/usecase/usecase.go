package usecase

import (
	"context"
	"dbms/internal/models"
	post "dbms/internal/post/repository"
	user "dbms/internal/user/repository"
	thread "dbms/internal/thread/repository"
	forum "dbms/internal/forum/repository"
)

type Usecase interface {
	GetFullPost(ctx context.Context, id int64, related []string) (models.FullPost, error)
	EditPost(ctx context.Context, id int64, updatePost models.UpdatePost) (models.Post, error)
}

type usecase struct {
	postRepository post.Repository
	userRepository  user.Repository
	threadRepository thread.Repository
	forumRepository forum.Repository
}

func stringInSlice(str string, list []string) bool {
    for _, s := range list {
        if s == str {
            return true
        }
    }
    return false
}

func (u *usecase) GetFullPost(ctx context.Context, id int64, related []string) (models.FullPost, error) {
	var fullPost models.FullPost

	post, err := u.postRepository.GetPostById(ctx, id)
	if err!= nil {
        return fullPost, err
    }
	
	fullPost.Post = post

	if stringInSlice("user", related) {
		author, err := u.userRepository.GetUserByNickname(ctx, post.Author)
		if err != nil {
			return fullPost, err
		}

		fullPost.Author = author
	}

	if stringInSlice("thread", related) {
		thread, err := u.threadRepository.GetThreadById(ctx, post.Thread)
		if err!= nil {
			return fullPost, err
		}

		fullPost.Thread = thread
	}

	if stringInSlice("forum", related) {
		forum, err := u.forumRepository.GetForumBySlug(ctx, post.Forum)
		if err!= nil {
			return fullPost, err
		}

		fullPost.Forum = forum
	}

	return fullPost, nil
}

func (u *usecase) EditPost(ctx context.Context, id int64, updatePost models.UpdatePost) (models.Post, error) {
	post, err := u.postRepository.EditPostById(ctx, id, updatePost)
	if err != nil {
        return models.Post{}, err
    }

	return post, nil
}

func NewUsecase (postRepository post.Repository) Usecase {
	return &usecase{postRepository: postRepository}
}
