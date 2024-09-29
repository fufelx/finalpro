// Пакет для работы с БД приложения GoNews.
package storage

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

// База данных.
type DB struct {
	pool *pgxpool.Pool
}

// Публикация, получаемая из RSS.
type Post struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	PubTime int64  `json:"pubtime"`
	Link    string `json:"link"`
}

type Pagi struct {
	Pages           int `json:"pages"`
	CurrentPage     int `json:"currentpage"`
	AmountOfElement int `json:"amountofelement"`
}

// Подключение к БД
func New() (*DB, error) {
	pool, err := pgxpool.Connect(context.Background(), "postgres://")
	if err != nil {
		return nil, err
	}
	db := DB{
		pool: pool,
	}
	return &db, nil
}

// Добавление новости в БД + в случае потоврения - пропуск
func (db *DB) StoreNews(news []Post) error {
	for _, post := range news {
		_, err := db.pool.Exec(context.Background(), `
		INSERT INTO news(title, content, pub_time, link)
		VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`,
			post.Title,
			post.Content,
			post.PubTime,
			post.Link,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// News возвращает последние новости из БД.
func (db *DB) News(n int) ([]Post, error) {
	if n == 0 {
		n = 10
	}
	rows, err := db.pool.Query(context.Background(), `
	SELECT id, title, content, pub_time, link FROM news
	ORDER BY pub_time DESC
	LIMIT $1
	`,
		n,
	)
	if err != nil {
		return nil, err
	}
	var news []Post
	for rows.Next() {
		var p Post
		err = rows.Scan(
			&p.Id,
			&p.Title,
			&p.Content,
			&p.PubTime,
			&p.Link,
		)
		if err != nil {
			return nil, err
		}
		news = append(news, p)
	}
	return news, rows.Err()
}

// News возвращает последние новости из БД.
func (db *DB) NewsByName(name string) ([]Post, error) {
	rows, err := db.pool.Query(context.Background(), `
		SELECT id, title, content, pub_time, link 
		FROM news 
		WHERE title ILIKE '%' || $1 || '%' 
		ORDER BY pub_time DESC`,
		name,
	)
	if err != nil {
		return nil, err
	}
	var news []Post
	for rows.Next() {
		var p Post
		err = rows.Scan(
			&p.Id,
			&p.Title,
			&p.Content,
			&p.PubTime,
			&p.Link,
		)
		if err != nil {
			return nil, err
		}
		news = append(news, p)
	}
	return news, rows.Err()
}

// News возвращает последние новости из БД.
func (db *DB) NewsById(id int) (Post, error) {
	rows, err := db.pool.Query(context.Background(), `
		SELECT id, title, content, pub_time, link 
		FROM news 
		WHERE id = $1`,
		id,
	)
	if err != nil {
		return Post{}, err
	}
	var news Post
	for rows.Next() {
		var p Post
		err = rows.Scan(
			&p.Id,
			&p.Title,
			&p.Content,
			&p.PubTime,
			&p.Link,
		)
		if err != nil {
			return Post{}, err
		}
		news = Post{p.Id,
			p.Title,
			p.Content,
			p.PubTime,
			p.Link}
	}
	return news, rows.Err()
}
