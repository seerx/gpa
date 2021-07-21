package models

import (
	"time"

	"github.com/seerx/gpa/examples/pratics/models/op"
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
	Data     []byte
}

type Teacher struct {
	ID       int    `gpa:"pk autoincr"`
	Name     string `gpa:"varchar(50)"`
	BirthDay time.Time
	Address  string `gpa:"varchar(200)"`
	Addr     op.Addr
}

// func (Student) Foo([]byte) error {
// 	return nil
// }

// func (s *Student) Read(buf []byte) error {
// 	return nil
// }

// func (s *Student) Write() (buf []byte, err error) {
// 	return nil, nil
// }
