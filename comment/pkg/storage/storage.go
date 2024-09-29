package storage

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"strings"
)

type Comment struct {
	Id      int
	Newsid  int
	Text    string
	Date    int64
	Parents int
	Allow   bool
}

type Store struct {
	ctx context.Context
	db  pgxpool.Pool
}

type Interface interface {
	NewComment(comment Comment) error // получение всех публикаций
}

func New(postgres_login string) (*Store, error) {
	var ctx context.Context = context.Background()
	db, err := pgxpool.Connect(ctx, postgres_login)
	if err != nil {
		return nil, err
	}
	result := Store{ctx: ctx, db: *db}
	return &result, nil
}

func (s *Store) NewComment(comment Comment) error {
	tx, err := s.db.Begin(s.ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(s.ctx)
	_, err = tx.Exec(s.ctx, `INSERT INTO comments(newsid, comment, dateunix, parents, allow) VALUES ($1,$2,$3,$4,$5)`, comment.Newsid, comment.Text, comment.Date, comment.Parents, WrongWord(comment.Text))
	if err != nil {
		return err
	} else {
		tx.Commit(s.ctx)
		return nil
	}
}

func (s *Store) AllComments(newsid int) ([]Comment, error) {
	rows, err := s.db.Query(s.ctx, `SELECT id, newsid, comment, dateunix, parents, allow FROM comments WHERE newsid = $1`, newsid)
	if err != nil {
		return nil, err
	} else {
		var comments []Comment
		var comment Comment
		for rows.Next() {
			var t Comment
			rows.Scan(&t.Id, &t.Newsid, &t.Text, &t.Date, &t.Parents, &t.Allow)
			comment = Comment{t.Id, t.Newsid, t.Text, t.Date, t.Parents, t.Allow}
			if rows.Err() != nil {
				return nil, err
			}
			comments = append(comments, comment)
		}
		return comments, nil
	}
}

func WrongWord(word string) bool {
	if strings.Contains(word, "qwerty") {
		return true
	}
	if strings.Contains(word, "йцукен") {
		return true
	}
	if strings.Contains(word, "zxvbnm") {
		return true
	}
	return false
}
