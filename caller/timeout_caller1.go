package main

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
	"xorm.io/xorm"
)

type User struct {
	Id      int64
	Name    string
	Salt    string
	Age     int
	Passwd  string    `xorm:"varchar(200)"`
	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
}

func CheckUserBool() bool {
	engine, err := xorm.NewEngine("mysql", "root:xxxxx@tcp(127.0.0.1:3306)/test?charset=utf8")

	if err != nil {
		log.Fatal(err)
	}

	/*
	   err = engine.Sync2(new(User))
	   if err != nil {
	     log.Fatal(err)
	   }

	         user := &User{Name: "panda", Age: 50}

	   affected, _ := engine.Insert(user)
	   fmt.Printf("%d records inserted, user.id:%d\n", affected, user.Id)

	*/
	user1 := &User{}
	has, _ := engine.ID(1).Get(user1)
	if has {
		fmt.Printf("user1:%v:%d\n", user1, user1.Age)
		if user1.Age > 60 {
			return true
		}
	}

	return false
}

func hardWork(job interface{}) error {
	for {
		ret := CheckUserBool()
		if ret {
			return nil
		}
		time.Sleep(2 * time.Second)
	}
	return nil
}

func requestBlockedWork(ctx context.Context, job interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*200)
	defer cancel()
	done := make(chan error, 1)
	panicChan := make(chan interface{}, 1)
	go func() {
		defer func() {
			if p := recover(); p != nil {
				panicChan <- p
			}
		}()
		done <- hardWork(job)
	}()
	select {
	case err := <-done:
		return err
	case p := <-panicChan:
		panic(p)
	case <-ctx.Done():
		return ctx.Err()
	}
}

func main() {
	now := time.Now()
	requestBlockedWork(context.Background(), "any")
	fmt.Println("elapsed:", time.Since(now))
}
