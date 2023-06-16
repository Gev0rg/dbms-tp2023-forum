package repository

import (
	"context"
	"dbms/internal/models"
	myErrors "dbms/internal/models/errors"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	CheckExistUserByNickname(ctx context.Context, nickname string) error
	CreateUser(ctx context.Context, createUser models.User) error
	UpdateUserByNickname(ctx context.Context, user models.User) (models.User, error)
	GetUserByNickname(ctx context.Context, nickname string) (models.User, error)
}

func NewUserMemoryRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

type repository struct {
	db *sqlx.DB
}

func (r repository) CheckExistUserByNickname(ctx context.Context, nickname string) error {
	var exists bool

	err := r.db.GetContext(ctx, &exists, "SELECT EXISTS(SELECT 1 FROM user u WHERE LOWER(u.nickname) = LOWER($1))", nickname)
	if err != nil {
		return err
	}

	if !exists {
		return myErrors.ErrUserNotFound
	}

	return nil
}

func (r repository) CreateUser(ctx context.Context, createUser models.User) error {
	_, err := r.db.NamedExecContext(ctx, `INSERT INTO user VALUES (:nickname, :fullname, :about, :email)`, createUser)
	return err
}

func (r repository) GetUserByNickname(ctx context.Context, nickname string) (models.User, error) {
	var user models.User

	err := r.db.GetContext(ctx, &user, `SELECT * FROM user WHERE nickname=$1`, nickname)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r repository) UpdateUserByNickname(ctx context.Context, user models.User) (models.User, error) {
	var updatedUser models.User

	err := r.db.GetContext(ctx, &updatedUser, `UPDATE user SET fullname=$1, about=$2, email=$3 WHERE nickname=$4 RETURNING *`,
		user.FullName, user.About, user.Email, user.Nickname)
	if err != nil {
		return models.User{}, err
	}

	return updatedUser, nil
}
