package core

import "context"

type User struct {
	ID     int64  `json:"id"`
	Login  string `json:"login"`
	Email  string `json:"email"`
	Admin  bool   `json:"admin"`
	Active bool   `json:"active"`
}

type UserStore interface {
	Find(context.Context, int64) (*User, error)
	Create(context.Context, *User) error
	Update(context.Context, *User) error
	Delete(context.Context, *User) error
}
