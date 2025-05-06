package pointers

import "fmt"

type User struct {
	email    string
	username string
	age      int
}

func (u User) Email() string {
	return u.email
}

func (u *User) UpdateEmail(email string) {
	u.email = email
}

func UpdateEmail(u *User, email string) {
	u.email = email
}

func Email(u User) string {
	return u.email
}

func HowPointerSampleOne() {
	user := User{
		email: "agg@foo.com",
	}

	user.UpdateEmail("ok.com")
	fmt.Println(user.Email())
}
