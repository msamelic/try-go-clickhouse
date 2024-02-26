package model

import "context"

type User struct {
	Id   uint32 `ch:"id"`
	Name string `ch:"name"`
	Age  uint8  `ch:"age"`
}

type Repository interface {
	ListUsers(context.Context) ([]User, error)
	AddUser(User) (User, error)
	GetUser(int) (User, error)
}
