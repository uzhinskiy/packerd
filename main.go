package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	//    "os/exec"
)

var (
	//NWorkers = flag.Int("n", 4, "The number of workers to start")
	HTTPAddr = flag.String("http", "127.0.0.1:8080", "Address to listen for HTTP requests on")
)

type vm_struct struct {
	Cloud  string           `json:"cloud,omitempty"`
	Region string           `json:"region,omitempty"`
	Vm     variables_struct `json:"vm,omitempty"`
}

type variables_struct struct {
	Templ_name string `json:"templ_name,omitempty"`
	Mem_vol    string `json:"mem_vol,omitempty"`
	Disk_size  string `json:"disk_size,omitempty"`
	Cpu_num    string `json:"cpu_num,omitempty"`
	Host_name  string `json:"host_name,omitempty"`
}

func packerCreate(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		{
			body, err := ioutil.ReadAll(req.Body)
			log.Println(string(body))

			var vm vm_struct
			var vars variables_struct
			//decoder := json.NewDecoder(req.Body)
			//err = decoder.Decode(&vars)

			err = json.Unmarshal(body, &vm)
			vars = vm.Vm

			defer req.Body.Close()

			if err != nil {
				log.Println(err)
			}

			if vars.Host_name != "" {
				fname := fmt.Sprintf("/tmp/variables_%s.json", vars.Templ_name)
				varsFile, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE, 0600)
				defer varsFile.Close()
				if err != nil {
					log.Println(err)
				}

				varsJson, _ := json.Marshal(vars)

				_, err = varsFile.WriteString(string(varsJson))
				if err != nil {
					log.Println(err)
				}
				rw.Header().Set("Content-Type", "application/json; charset=utf-8")
				rw.Header().Set("Server", "packed/0.1")
				rw.WriteHeader(http.StatusCreated)
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

func packerStatus(rw http.ResponseWriter, req *http.Request) {
	log.Println(req.URL)
	switch req.Method {
	case "GET":
		{
			req.ParseForm()
			// logic part of log in
			fmt.Println("fname:", req.Form)
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
	flag.Parse()
	http.HandleFunc("/status", packerStatus)
	http.HandleFunc("/create", packerCreate)
	log.Println("HTTP server listening on", *HTTPAddr)
	err := http.ListenAndServe(*HTTPAddr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
