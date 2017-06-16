package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	//    "os"
	//    "os/exec"
)

type test_struct struct {
	Test string
}

func packerJson(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(body))
	var t test_struct
	err = json.Unmarshal(body, &t)
	if err != nil {
		log.Println(err)
	}
	log.Println(t.Test)
	fmt.Fprintf(rw, "OK\n")
}

func packerPost(rw http.ResponseWriter, req *http.Request) {
	log.Println(req)
	if req.Method == "GET" {
		fmt.Fprintf(rw, "wrong method")
	} else {
		log.Println(req)
		req.ParseForm()
		// logic part of log in
		fmt.Println("fname:", req.PostFormValue("fname"))
		fmt.Fprintf(rw, "OK\n")
	}
}

func main() {
	http.HandleFunc("/post", packerPost)
	http.HandleFunc("/json", packerJson)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
