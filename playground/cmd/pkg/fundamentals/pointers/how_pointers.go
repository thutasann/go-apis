/*
# When we use Pointers

1. When we need to update state

2. When we want to optimize the memory for large objects that are getting called A LOT
*/
package pointers

import (
	"fmt"
	"unsafe"
)

type User struct {
	email    string
	username string
	age      int
	file     []byte // ?? dont know how large ??
}

func GetUser() (*User, error) {
	return nil, fmt.Errorf("foo")
}

// x amount of bytes => sizeOf(user)
func (u User) Email() string {
	return u.email
}

func (u User) UserName() string {
	return u.username
}

func (u *User) UpdateEmail(email string) {
	u.email = email
}

func (u *User) UpdateName(name string) {
	u.username = name
}

func UpdateEmail(u *User, email string) {
	u.email = email
}

func Email(u User) string {
	return u.email
}

// When & How to use Pointers
func HowPointerSampleOne() {
	user := User{
		email:    "agg@foo.com",
		username: "agg",
		age:      20,
	}

	user_size := unsafe.Sizeof(user)
	fmt.Println("user_size :>>", user_size)

	user.UpdateEmail("ok.com")
	user.UpdateName("agg updated")
	fmt.Println(user.Email())
	fmt.Println(user.UserName())
}
