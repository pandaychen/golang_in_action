package test

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"testing"
	"time"
	"xorm.io/core"
	"xorm.io/xorm"
)

type user struct {
	ID        uint64 `xorm:"autoincr"`
	FirstName string
	LastName  string
	Extra1    string
	Extra2    string
	Version   uint64    `xorm:"version"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	DeletedAt time.Time `xorm:"deleted"`
}

func (user) TableName() string {
	return "users"
}

func insert(db xorm.Interface, us users) error {
	_, err := db.Insert(&us)
	return err
}

type users []user

func do_insert(b *testing.B, insertCount int, fn func(xorm.Interface, users) error) {
	us := make(users, 0, insertCount)
	for i := 0; i < insertCount; i++ {
		us = append(us, user{
			FirstName: fmt.Sprintf("first_%v", i),
			LastName:  fmt.Sprintf("last_%v", i),
			Extra1:    fmt.Sprintf("extra1_%v", i),
			Extra2:    fmt.Sprintf("extra2_%v", i),
		})
	}

	engine, err := xorm.NewEngine("mysql", "user:xxxxx@tcp(127.0.0.1:3306)/test_db?charset=utf8")
	if err != nil {
		panic(err)
	}
	defer engine.Close()
	engine.SetMapper(core.GonicMapper{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		_, err := engine.Transaction(func(tx *xorm.Session) (interface{}, error) {
			err := fn(tx, us)
			return nil, err
		})
		if err != nil {
			panic(err)
		}
		b.StopTimer()
		engine.Unscoped().Where("id > ?", 0).Delete(user{})
	}
}

func BenchmarkInsert1(b *testing.B) {
	do_insert(b, 1, insert)
}

func BenchmarkInsert10(b *testing.B) {
	do_insert(b, 10, insert)
}

func BenchmarkInsert100(b *testing.B) {
	do_insert(b, 100, insert)
}

func BenchmarkInsert1000(b *testing.B) {
	do_insert(b, 1000, insert)
}
