// Пакет для работы с БД приложения GoNews.

package storage

import (
	"log"
	"news/pkg/storage"
	"testing"
)

func TestDB_StoreNews(t *testing.T) {
	db, _ := storage.New()
	type args struct {
		news []storage.Post
	}
	tests := []struct {
		name    string
		db      *storage.DB
		args    args
		wantErr bool
	}{
		{
			name: "Valid DB",
			db:   db,
			args: args{
				news: []storage.Post{
					{Title: "Test", Content: "Test", PubTime: 0, Link: "Test"},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.db.StoreNews(tt.args.news); (err != nil) != tt.wantErr {
				t.Errorf("DB.StoreNews() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDB_News(t *testing.T) {
	db, _ := storage.New()
	type args struct {
		n int
	}
	tests := []struct {
		name    string
		db      *storage.DB
		args    args
		want    []Post
		wantErr bool
	}{
		{
			name:    "Valid DB",
			db:      db,
			args:    args{5},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.db.News(tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("DB.News() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else {
				log.Println(got)
			}
			/*
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("DB.News() = %v, want %v", got, tt.want)
				}
			*/
		})
	}
}
