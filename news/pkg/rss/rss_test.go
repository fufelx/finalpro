package rss

import (
	"log"
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		rssURL string
	}
	tests := []struct {
		name string
		args args
		//want    [...]storage.Post
		wantErr bool
	}{
		{
			name: "Valid RSS URL",
			args: args{"https://rsshub.app/telegram/channel/mash"},
			//want:    []storage.Post{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.rssURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else {
				log.Println(got)
			}
			/*
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Parse() = %v, want %v", got, tt.want)
				}
			*/

		})
	}
}
