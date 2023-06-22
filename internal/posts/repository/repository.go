package repository

import (
	"context"
	"database/sql"
	"fmt"
	myerr "forum/internal/error"
	"forum/internal/models"
	"forum/internal/posts"
	"log"
	"regexp"
	"strings"

	"github.com/lib/pq"
)

type PostRepository struct {
	db     *sql.DB
	logger *log.Logger
}

func NewPostRepository(db *sql.DB) posts.PostRepository {
	return &PostRepository{
		db:     db,
		logger: log.Default(),
	}
}

func (pr *PostRepository) SelectFormSlugByThread(slug string, id int64) (string, int64, error) {
	queryStr := "SELECT forum, id FROM threads WHERE 0 = $1 AND slug = $2 OR $2 = '' AND id = $1"
	row := pr.db.QueryRow(queryStr, id, slug)
	var forumSlug string
	var threadId int64
	err := row.Scan(&forumSlug, &threadId)
	if err != nil {
		res, _ := regexp.Match(".*no rows in result set.*", []byte(err.Error()))
		if res {
			return "", 0, myerr.ThreadNotExists
		}
		pr.logger.Println(err.Error())
		return "", 0, myerr.InternalDbError
	}
	return forumSlug, threadId, nil
}

func (pr *PostRepository) CheckNickname(nickname string) (string, error) {
	row := pr.db.QueryRow("SELECT nickname FROM users WHERE nickname = $1;", nickname)
	err := row.Scan(&nickname)
	if err != nil {
		res, _ := regexp.Match(".*no rows in result set.*", []byte(err.Error()))
		if res {
			return "", myerr.UserNotExist
		}
		pr.logger.Println(err.Error())
		return "", myerr.InternalDbError
	}
	return nickname, nil
}

func (pr *PostRepository) CheckParent(threadId int64, parent int64) ([]int64, error) {
	path := make([]int64, 0)
	row := pr.db.QueryRow("SELECT array_append(path, id) FROM posts WHERE id = $1 AND thread = $2;", parent, threadId)
	err := row.Scan(pq.Array(&path))
	if err != nil {
		res, _ := regexp.Match(".*no rows in result set.*", []byte(err.Error()))
		if res {
			return nil, myerr.ParentNotExist
		}
		pr.logger.Println(err.Error())
		return nil, myerr.InternalDbError
	}
	return path, nil
}

func (pr *PostRepository) CreatePosts(inputPost []*models.PostInput, dt string, forumSlug string, threadId int64) ([]*models.Post, error) {
	queryStr := "INSERT INTO posts (message, forum, thread, created, author, parent, path) VALUES"
	args := make([]interface{}, 0)
	var err1, err2 error = nil, nil
	for ind, ip := range inputPost {
		path := make([]int64, 0)
		ip.Author, err1 = pr.CheckNickname(ip.Author)
		if ip.Parent != 0 {
			path, err2 = pr.CheckParent(threadId, ip.Parent)
		} else {
			err2 = nil
		}

		if err1 != nil {
			return nil, err1
		}

		if err2 != nil {
			return nil, err2
		}

		queryStr += fmt.Sprintf(
			" ($%d, $%d, $%d, $%d, $%d, $%d, CAST($%d AS BIGINT ARRAY)),",
			ind*7+1, ind*7+2, ind*7+3, ind*7+4, ind*7+5, ind*7+6, ind*7+7)
		args = append(args, ip.Message, forumSlug, threadId, dt, ip.Author, ip.Parent, pq.Array(path))
	}

	tx, err := pr.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		pr.logger.Println(err.Error())
		return nil, myerr.InternalDbError
	}

	queryStr = strings.TrimSuffix(queryStr, ",")
	queryStr += " RETURNING id, message, forum, thread, created, author, parent, isEdited;"
	rows, err := tx.Query(queryStr, args...)
	if err != nil {
		pr.logger.Println("before scan:", err.Error())
		return nil, myerr.InternalDbError
	}
	defer rows.Close()

	posts := make([]*models.Post, 0)
	for rows.Next() {
		post := &models.Post{}
		err = rows.Scan(&post.Id, &post.Message, &post.Forum, &post.Thread, &post.Created, &post.Author, &post.Parent, &post.IsEdited)
		if err != nil {
			rollbackError := tx.Rollback()
			if rollbackError != nil {
				return nil, myerr.RollbackError
			}

			pr.logger.Println(err.Error())
			return nil, myerr.InternalDbError
		}

		posts = append(posts, post)
	}

	err = tx.Commit()
	if err != nil {
		return nil, myerr.CommitError
	}

	return posts, nil
}

func (pr *PostRepository) CreatePost(inputPost *models.PostInput, dt string, forumSlug string, threadId int64) (*models.Post, error) {
	tx, err := pr.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		pr.logger.Println(err.Error())
		return nil, myerr.InternalDbError
	}

	var row *sql.Row
	queryStr := `
		INSERT INTO posts (message, forum, thread, created, author%s) 
		VALUES ($1, $2, $3, $4,
			COALESCE((SELECT nickname FROM users WHERE nickname = $5), $5)
			%s
		)
		RETURNING id, message, forum, thread, created, author, parent, isEdited;`
	if inputPost.Parent == 0 {
		queryStr = fmt.Sprintf(queryStr, "", "")
		row = tx.QueryRow(queryStr, inputPost.Message, forumSlug, threadId, dt, inputPost.Author)
	} else {
		queryStr = fmt.Sprintf(
			queryStr,
			", parent",
			", (SELECT (CASE WHEN EXISTS (SELECT id FROM posts WHERE id = $6 AND thread = $3) THEN $6 ELSE 0 END))")
		row = tx.QueryRow(queryStr, inputPost.Message, forumSlug, threadId, dt, inputPost.Author, inputPost.Parent)
	}

	post := &models.Post{}
	err = row.Scan(&post.Id, &post.Message, &post.Forum, &post.Thread, &post.Created, &post.Author, &post.Parent, &post.IsEdited)
	if err != nil {
		rollbackError := tx.Rollback()
		if rollbackError != nil {
			return nil, myerr.RollbackError
		}

		res, _ := regexp.Match(".*\"parent\" violates.*", []byte(err.Error()))
		if res {
			return nil, myerr.ParentNotExist
		}

		res, _ = regexp.Match(".*posts_author_fkey.*", []byte(err.Error()))
		if res {
			return nil, myerr.UserNotExist
		}

		pr.logger.Println(err.Error())
		return nil, myerr.InternalDbError
	}

	err = tx.Commit()
	if err != nil {
		return nil, myerr.CommitError
	}
	return post, nil
}

func (pr *PostRepository) SelectThread(id int64, slug string) (int64, error) {
	row := pr.db.QueryRow(
		"SELECT id from threads WHERE 0 = $1 AND slug = $2 OR $2 = '' AND id = $1",
		id, slug)
	err := row.Scan(&id)
	if err != nil {
		res, _ := regexp.Match(".*no rows in result set.*", []byte(err.Error()))
		if res {
			return 0, myerr.ThreadNotExists
		}
		return 0, myerr.InternalDbError
	}
	return id, nil
}

func (pr *PostRepository) SelectThreadsBySort(tq *models.ThreadsQuery) ([]*models.Post, error) {
	var queryStr string
	var counter uint = 2
	var nums []interface{}
	var args []interface{}
	if tq.Sort == "flat" {
		queryStr = `SELECT id, message, forum, thread, created, author, parent, isEdited 
					FROM posts WHERE thread = $1 `
		args = append(args, tq.ThreadId)
		if tq.Since != 0 {
			queryStr = queryStr + "AND id %s $%d "
			nums = append(nums, tq.Sign, counter)
			counter = counter + 1
			args = append(args, tq.Since)
		}
		queryStr = queryStr + "ORDER BY created %s, id %s LIMIT $%d"
		nums = append(nums, tq.Sorting, tq.Sorting, counter)
		args = append(args, tq.Limit)
		queryStr = fmt.Sprintf(queryStr, nums...)
	} else if tq.Sort == "tree" {
		queryStr = `SELECT id, message, forum, thread, created, author, parent, isEdited 
					FROM posts WHERE thread = $1 `
		args = append(args, tq.ThreadId)
		if tq.Since != 0 {
			queryStr = queryStr + "AND array_append(path, id) %s (SELECT array_append(path, id) FROM posts WHERE id = $%d) "
			nums = append(nums, tq.Sign, counter)
			counter = counter + 1
			args = append(args, tq.Since)
		}
		queryStr = queryStr + "ORDER BY array_append(path, id) %s LIMIT $%d"
		nums = append(nums, tq.Sorting, counter)
		args = append(args, tq.Limit)
		queryStr = fmt.Sprintf(queryStr, nums...)
	} else if tq.Sort == "parent_tree" {
		queryStr = `SELECT t.id, t.message, t.forum, t.thread, t.created, t.author, t.parent, t.isEdited 
					FROM (SELECT *, CASE WHEN cardinality(path) = 0 THEN id ELSE path[1] END as rooot FROM posts) as t
					WHERE t.rooot IN (
						SELECT id FROM posts
						WHERE thread = $1 AND parent = 0 %s
						ORDER BY id %s
						LIMIT $%d
					) AND thread = $1
					ORDER BY t.rooot %s, array_append(path, id)`
		args = append(args, tq.ThreadId)
		s1 := ""
		if tq.Since != 0 {
			s1 = "AND id %s (SELECT CASE WHEN cardinality(path) = 0 THEN id ELSE path[1] END FROM posts WHERE id = $%d)"
			s1 = fmt.Sprintf(s1, tq.Sign, counter)
			counter = counter + 1
			args = append(args, tq.Since)
		}
		args = append(args, tq.Limit)
		queryStr = fmt.Sprintf(queryStr, s1, tq.Sorting, counter, tq.Sorting)
	}

	rows, err := pr.db.Query(queryStr, args...)
	if err != nil {
		pr.logger.Println(err.Error())
		return nil, myerr.InternalDbError
	}

	posts := make([]*models.Post, 0)
	for rows.Next() {
		post := &models.Post{}
		err = rows.Scan(&post.Id, &post.Message, &post.Forum, &post.Thread, &post.Created, &post.Author, &post.Parent, &post.IsEdited)
		if err != nil {
			pr.logger.Println(err.Error())
			return nil, myerr.InternalDbError
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (pr *PostRepository) SelectPost(id int64) (*models.Post, error) {
	post := &models.Post{}
	row := pr.db.QueryRow(
		"SELECT id, parent, author, message, isEdited, forum, thread, created FROM posts WHERE id = $1;",
		id)
	err := row.Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created)
	if err != nil {
		res, _ := regexp.Match(".*no rows in result set.*", []byte(err.Error()))
		if res {
			return nil, myerr.PostNotExist
		}
		return nil, myerr.InternalDbError
	}
	return post, nil
}

func (pr *PostRepository) SelectUser(nickname string) (*models.User, error) {
	user := &models.User{}
	row := pr.db.QueryRow(
		"SELECT nickname, fullname, email, about FROM users WHERE nickname = $1;",
		nickname)
	err := row.Scan(&user.Nickname, &user.Fullname, &user.Email, &user.About)
	if err != nil {
		res, _ := regexp.Match(".*no rows in result set.*", []byte(err.Error()))
		if res {
			return nil, myerr.UserNotExist
		}
		return nil, myerr.InternalDbError
	}
	return user, nil
}

func (pr *PostRepository) SelectThreadById(id int64) (*models.Thread, error) {
	thread := &models.Thread{}
	row := pr.db.QueryRow(
		"SELECT id, slug, title, author, forum, message, votes, created FROM threads WHERE id = $1;",
		id)
	err := row.Scan(&thread.Id, &thread.Slug, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Created)
	if err != nil {
		res, _ := regexp.Match(".*no rows in result set.*", []byte(err.Error()))
		if res {
			return nil, myerr.ThreadNotExists
		}
		return nil, myerr.InternalDbError
	}
	return thread, nil
}

func (pr *PostRepository) SelectForum(slug string) (*models.Forum, error) {
	forum := &models.Forum{}
	row := pr.db.QueryRow(
		"SELECT slug, title, author, posts, threads FROM forum WHERE slug = $1;",
		slug)
	err := row.Scan(&forum.Slug, &forum.Title, &forum.User, &forum.Posts, &forum.Threads)
	if err != nil {
		res, _ := regexp.Match(".*no rows in result set.*", []byte(err.Error()))
		if res {
			return nil, myerr.ForumNotExist
		}
		return nil, myerr.InternalDbError
	}
	return forum, nil
}

func (pr *PostRepository) UpdatePost(postupdate *models.PostUpdate) (*models.Post, error) {
	tx, err := pr.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		pr.logger.Println(err.Error())
		return nil, myerr.InternalDbError
	}

	post := &models.Post{}
	row := tx.QueryRow(
		`UPDATE posts SET 
		 	message = CASE WHEN $2 = '' THEN message ELSE $2 END, 
			isEdited = CASE WHEN $2 = '' THEN isEdited ELSE CASE WHEN message = $2 THEN isEdited ELSE $3 END END 
		 WHERE id = $1
		 RETURNING id, parent, author, message, isEdited, forum, thread, created;`,
		postupdate.Id, postupdate.Message, true)
	err = row.Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created)
	if err != nil {
		rollbackError := tx.Rollback()
		if rollbackError != nil {
			return nil, myerr.RollbackError
		}
		res, _ := regexp.Match(".*no rows in result set.*", []byte(err.Error()))
		if res {
			return nil, myerr.PostNotExist
		}

		pr.logger.Println(err.Error())
		return nil, myerr.InternalDbError
	}

	err = tx.Commit()
	if err != nil {
		return nil, myerr.CommitError
	}
	return post, nil
}
