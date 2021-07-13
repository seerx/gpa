package models

import (
	"time"
)

type User struct {
	ID       uint64 `gpa:"pk not-null"`
	Name     string `gpa:"varchar(202)"`
	Age      int    `gpa:"index(age0)"`
	Gender   bool   `gpa:"index(age0)"`
	OkaUsURL bool   `gpa:"default(true)"`
	Cret     time.Time
	// What     ttt.Abc
	URL string
}

type Student struct {
	ID       uint64 `gpa:"pk not-null"`
	Name     string `gpa:"varchar(50)"`
	ClassID  uint64
	BirthDay time.Time
	Address  string `gpa:"varchar(200)"`
}
