// redis.go
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/garyburd/redigo/redis"
)

func RedisConnect() redis.Conn {
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		log.Fatalln(err)
	}
	return c
}

func RedisSet(w WorkEntry) {
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

func RedisGet(id string) WorkEntry {
	c := RedisConnect()
	defer c.Close()
	var w WorkEntry

	// Save JSON to Redis
	reply, err := c.Do("GET", "jobs:"+id)
	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal(reply.([]byte), &w)
	if err != nil {
		log.Println(err)
	}
	return w
}
