package main

import (
	"github.com/chrisseto/evil"
	"github.com/chrisseto/evil/channel"
)

type PostsView struct {
}

func (v *PostsView) OnMount(*evil.Session) error {
	return nil
}

func (v *PostsView) ToArgs(*evil.Session) (interface{}, error) {
	return map[string]interface{}{
		"Posts": []Post{
			{ID: "1", Body: "Hello World", Username: "chris", Likes: 2},
			{ID: "2", Body: "Good Bye World", Username: "jill", Likes: 0},
		},
	}, nil
}

func (v *PostsView) HandleEvent(*evil.Session, *channel.Event) error {
	return nil
}
