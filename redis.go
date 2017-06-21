// redis.go
package main

import (
	//"encoding/json"
	"fmt"
	"log"

	"github.com/garyburd/redigo/redis"
)

type Storage struct {
	db redis.Conn
}

// http://www.alexedwards.net/blog/organising-database-access
// https://www.reddit.com/r/golang/comments/38hkor/go_best_practice_for_accessing_database_in/
// https://stackoverflow.com/questions/41257847/how-to-create-singleton-db-class-in-golang
func (s Storage) New() error {
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return err
	}
	fmt.Println("c", c)
	s = Storage{db: c}
	fmt.Println("s", s)
	return nil
}

func (s Storage) Close() error {
	err := s.db.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s Storage) Put(uid string, ss string) {
	// Save JSON to Redis
	reply, err := s.db.Do("SET", "jobs:"+uid, ss)
	fmt.Println(reply)
	if err != nil {
		log.Println(err)
	}
}

/*
func RedisConnect() redis.Conn {
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		log.Println(err)
	}
	return c
}

func (s *Storage) Put(w WorkRequest) {
	c := RedisConnect()
	defer c.Close()

	// Marshal Post to JSON blob
	b, err := json.Marshal(w)

	// Save JSON to Redis
	reply, err := c.Do("SET", "jobs:"+w.UID, b)
	fmt.Println(reply)
	if err != nil {
		log.Println(err)
	}
}


*/
