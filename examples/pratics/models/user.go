package models

import (
	"time"
)

type User struct {
	ID       uint64 `gpa:"pk"`
	Name     string `gpa:"varchar(202)"`
	Age      int    `gpa:"index(age0)"`
	Gender   bool   `gpa:"index(age0)"`
	OkaUsURL bool   `gpa:"default(true)"`
	Cret     time.Time
	// What     ttt.Abc
	URL string
}
