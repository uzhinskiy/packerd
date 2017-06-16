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

type variables_struct struct {
	xenserver_user     string
	xenserver_password string
	templ_name         string
	mem_vol            int
	disk_size1         int
	cpu_num            string
	host_name          string
}

func packerJson(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(body))
	var vars variables_struct
	err = json.Unmarshal(body, &vars)
	if err != nil {
		log.Println(err)
	}
	log.Println(vars.host_name)
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
