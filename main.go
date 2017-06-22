package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
)

var (
	NWorkers = flag.Int("n", 4, "The number of workers to start")
	HTTPAddr = flag.String("http", "127.0.0.1:8080", "Address to listen for HTTP requests on")
)

type WorkRequest struct {
	Platform string
	Region   string
	UID      string
}

type WorkEntry struct {
	Status string
	UID    string
}

type vm_struct struct {
	Platform string           `json:"platform,omitempty"`
	Region   string           `json:"region,omitempty"`
	UID      string           `json:"UID,omitempty"`
	Vars     variables_struct `json:"vars,omitempty"`
}

type variables_struct struct {
	Templ_name string `json:"templ_name,omitempty"`
	Mem_vol    string `json:"mem_vol,omitempty"`
	Disk_size  string `json:"disk_size,omitempty"`
	Cpu_num    string `json:"cpu_num,omitempty"`
	Host_name  string `json:"host_name,omitempty"`
}

// Буфферизиованный канал через который передаются задания.
var WorkQueue = make(chan WorkRequest, 100)
var ResponseQueue = make(chan WorkEntry, 100)

func packerCreate(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		{
			var vm vm_struct
			var vars variables_struct

			body, err := ioutil.ReadAll(req.Body)

			err = json.Unmarshal(body, &vm)
			if err != nil {
				log.Println(err)
				return
			}

			vars = vm.Vars
			varsJson, _ := json.Marshal(vars)

			defer req.Body.Close()

			fname := fmt.Sprintf("/tmp/%s.json", vm.UID)
			varsFile, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE, 0600)
			defer varsFile.Close()
			if err != nil {
				log.Println(err)
				return
			}

			_, err = varsFile.WriteString(string(varsJson))
			if err != nil {
				log.Println(err)
				return
			}

			// создание структуры для worker-а
			work := WorkRequest{UID: vm.UID, Region: vm.Region, Platform: vm.Platform}
			WorkQueue <- work
			log.Println("Work request queued")

			rw.Header().Set("Content-Type", "application/json; charset=utf-8")
			rw.Header().Set("Server", "packerd/0.1")
			rw.WriteHeader(http.StatusCreated)
			fmt.Fprintln(rw, "{\"status\":\"ok\", \"UID\":\""+vm.UID+"\"}")

		}
	default:
		{
			rw.Header().Set("Content-Type", "text/html; charset=utf-8")
			rw.Header().Set("Server", "packerd/0.1")
			rw.Header().Set("Allow", "POST")
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	}
}

func packerStatus(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		{
			var resp WorkEntry

			req.ParseForm()
			uid := path.Base(req.URL.Path)

			select {
			case resp := <-ResponseQueue:
				log.Println(resp)
			default:
				resp = RedisGet(uid)
			}

			rw.Header().Set("Content-Type", "application/json; charset=utf-8")
			rw.Header().Set("Access-Control-Allow-Origin", "*")
			rw.Header().Set("Server", "packerd/0.1")
			json, _ := json.Marshal(resp)
			fmt.Fprintln(rw, string(json))
		}
	default:
		{
			rw.Header().Set("Content-Type", "text/html; charset=utf-8")
			rw.Header().Set("Server", "packerd/0.1")
			rw.WriteHeader(http.StatusNotImplemented)
		}
	}
}

func main() {
	flag.Parse()

	for w := 1; w <= *NWorkers; w++ {
		go worker(w, WorkQueue, ResponseQueue)
	}

	http.HandleFunc("/status/", packerStatus)
	http.HandleFunc("/create", packerCreate)
	log.Println("HTTP server listening on", *HTTPAddr)
	err := http.ListenAndServe(*HTTPAddr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
