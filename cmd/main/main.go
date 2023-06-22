package main

import (
	forumDelivery "forum/internal/forum/delivery"
	"forum/internal/middleware"
	postDelivery "forum/internal/posts/delivery"
	serviceDelivery "forum/internal/service/delivery"
	threadDelivery "forum/internal/threads/delivery"
	userDelivery "forum/internal/user/delivery"
	voteDelivery "forum/internal/votes/delivery"

	forumUsecase "forum/internal/forum/usecase"
	postUsecase "forum/internal/posts/usecase"
	serviceUsecase "forum/internal/service/usecase"
	threadUsecase "forum/internal/threads/usecase"
	userUsecase "forum/internal/user/usecase"
	voteUsecase "forum/internal/votes/usecase"

	forumRepository "forum/internal/forum/repository"
	postRepository "forum/internal/posts/repository"
	serviceRepository "forum/internal/service/repository"
	threadRepository "forum/internal/threads/repository"
	userRepository "forum/internal/user/repository"
	voteRepository "forum/internal/votes/repository"

	"fmt"
	"forum/db"
	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/stdlib"
	"log"
	"net/http"
)

func main() {
	dbConnStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", "gev", "gev", "127.0.0.1", "5432", "gev")
	db, err := db.NewDatabase(dbConnStr)
	if err != nil {
		return
	}
	defer db.Close()

	forumRepository := forumRepository.NewForumRepository(db)
	userRepository := userRepository.NewUserRepository(db)
	threadRepository := threadRepository.NewThreadRepository(db)
	postRepository := postRepository.NewPostRepository(db)
	voteRepository := voteRepository.NewVoteRepository(db)
	serviceRepository := serviceRepository.NewServiceRepository(db)

	forumUsecase := forumUsecase.NewForumUsecase(forumRepository)
	userUsecase := userUsecase.NewUserUsecase(userRepository)
	threadUsecase := threadUsecase.NewThreadUsecase(threadRepository)
	postUsecase := postUsecase.NewPostUsecase(postRepository)
	voteUsecase := voteUsecase.NewVoteUsecase(voteRepository)
	serviceUsecase := serviceUsecase.NewServiceUsecase(serviceRepository)

	forumDelivery := forumDelivery.NewForumDelivery(forumUsecase)
	userDelivery := userDelivery.NewUserDelivery(userUsecase)
	threadDelivery := threadDelivery.NewForumDelivery(threadUsecase)
	postDelivery := postDelivery.NewPostDelivery(postUsecase)
	voteDelivery := voteDelivery.NewVoteDelivery(voteUsecase)
	serviceDelivery := serviceDelivery.NewServiceDelivery(serviceUsecase)

	r := mux.NewRouter()
	r = r.PathPrefix("/api").Subrouter()
	r.Use(middleware.ContentTypeMiddleware)

	forumDelivery.Routing(r)
	userDelivery.Routing(r)
	threadDelivery.Routing(r)
	postDelivery.Routing(r)
	voteDelivery.Routing(r)
	serviceDelivery.Routing(r)

	port := 5000
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), r)
	if err != nil {
		log.Default().Printf("start serving ::%d\n", port)
	} else {
		log.Default().Fatalf("http serve error %v", err)
	}
}
