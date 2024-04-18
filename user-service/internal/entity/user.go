package entity

import "time"

type User struct {
	GUID      string
	FirstName string
	LastName  string
	Username  string
	Email     string
	Password  string
	Bio       string
	Website   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
