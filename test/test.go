package main

import "fmt"

type User struct {
	name string
	//age int
}

func (u *User) Set(name string) {
	u.name = name
}

func (u *User) Get() string {
	return u.name
}

func NewUser() User {
	return User{}
}

func main() {
	user := NewUser()
	user.Set("tom")
	fmt.Println(user.Get())
}
