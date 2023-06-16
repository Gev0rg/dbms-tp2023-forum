package main

import (
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"os"
	"dbms/internal/middleware"
	"dbms/configs"

	// forumRepository "dbms/internal/forum/repository"
	// postRepository "dbms/internal/post/repository"
	// serviceRepository "dbms/internal/service/repository"
	// threadRepository "dbms/internal/thread/repository"
	userRepository "dbms/internal/user/repository"

	// forumUsecase "dbms/internal/forum/usecase"
	// postUsecase "dbms/internal/post/usecase"
	// serviceUsecase "dbms/internal/service/usecase"
	// threadUsecase "dbms/internal/thread/usecase"
	userUsecase "dbms/internal/user/usecase"

	// forumHandler "dbms/internal/forum/delivery/http"
	// postHandler "dbms/internal/post/delivery/http"
	// serviceHandler "dbms/internal/service/delivery/http"
	// threadHandler "dbms/internal/thread/delivery/http"
	userHandler "dbms/internal/user/delivery/http"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

func init() {
	envPath := ".env"

	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatal("No .env file found")
	}
}

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	yamlPath, exists := os.LookupEnv("YAML_PATH")
	if !exists {
		log.Fatal("Yaml path not found")
	}

	yamlFile, err := os.ReadFile(yamlPath)
	if err != nil {
		log.Fatal(err)
	}

	var config configs.Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sqlx.Open(config.Postgres.DB, config.Postgres.ConnectionToDB)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)

	e := echo.New()
	e.Use(middleware.HandlerMiddleware)

	// forumRepository := forumRepository.NewForumMemoryRepository(db)
	// postRepository := postRepository.NewPostMemoryRepository(db)
	// serviceRepository := serviceRepository.NewServiceMemoryRepository(db)
	// threadRepository := threadRepository.NewThreadMemoryRepository(db)
	userRepository := userRepository.NewUserMemoryRepository(db)

	// forumUsecase := forumUsecase.NewForumUsecase(forumRepository, userRepository)
	// postUsecase := postUsecase.NewPostUsecase(postRepository)
	// serviceUsecase := serviceUsecase.NewServiceUsecase(serviceRepository)
	// threadUsecase := threadUsecase.NewThreadUsecase(threadRepository)
	userUsecase := userUsecase.NewUserUsecase(userRepository)

	// forumHandler.NewForumHandler(e, forumUsecase)
	// postHandler.NewPostHandler(e, postUsecase)
	// serviceHandler.NewServiceHandler(e, serviceUsecase)
	// threadHandler.NewThreadHandler(e, threadUsecase)
	userHandler.NewHandler(e, userUsecase)

	e.Logger.Fatal(e.Start(config.Server.Port))
}
