package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	//    "os/exec"
)

type variables_struct struct {
	Xenserver_user     string `json:"xenserver_user,omitempty"`
	Xenserver_password string `json:"xenserver_password,omitempty"`
	Templ_name         string `json:"templ_name,omitempty"`
	Mem_vol            int    `json:"mem_vol,omitempty"`
	Disk_size1         int    `json:"disk_size1,omitempty"`
	Cpu_num            string `json:"cpu_num,omitempty"`
	Host_name          string `json:"host_name,omitempty"`
}

func packerJson1(rw http.ResponseWriter, req *http.Request) {
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
	log.Println(vars.Host_name)
	fmt.Fprintf(rw, "OK\n")
}

func packerJson(rw http.ResponseWriter, req *http.Request) {
	log.Println(req)
	switch req.Method {
	case "POST":
		{
			body, err := ioutil.ReadAll(req.Body)
			log.Println(string(body))

			var vars variables_struct
			//decoder := json.NewDecoder(req.Body)
			//err = decoder.Decode(&vars)

			err = json.Unmarshal(body, &vars)

			defer req.Body.Close()

			if err != nil {
				log.Println(err)
			}

			if vars.Host_name != "" {
				fname := fmt.Sprintf("/tmp/variables_%s.json", vars.Templ_name)
				vars_file, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE, 0600)
				defer vars_file.Close()
				if err != nil {
					log.Println(err)
				}

				_, err = vars_file.WriteString(string(body))
				if err != nil {
					log.Println(err)
				}
				rw.Header().Set("Content-Type", "application/json; charset=utf-8")
				rw.Header().Set("Server", "packed/0.1")
				fmt.Fprint(rw, "{\"status\":\"ok\"}")

			}
		}
	default:
		{
			rw.Header().Set("Content-Type", "text/html; charset=utf-8")
			rw.Header().Set("Server", "packed/0.1")
			rw.WriteHeader(http.StatusNotImplemented)
		}
	}
}

func packerPost(rw http.ResponseWriter, req *http.Request) {
	log.Println(req)
	switch req.Method {
	case "POST":
		{
			req.ParseForm()
			// logic part of log in
			fmt.Println("fname:", req.PostForm["fname"])
			rw.Header().Set("Content-Type", "application/json; charset=utf-8")
			rw.Header().Set("Server", "packed/0.1")
			fmt.Fprint(rw, "{\"status\":\"ok\"}")
		}
	default:
		{
			rw.Header().Set("Content-Type", "text/html; charset=utf-8")
			rw.Header().Set("Server", "packed/0.1")
			rw.WriteHeader(http.StatusNotImplemented)
		}
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
