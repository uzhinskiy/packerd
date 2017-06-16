package main

import (
    "encoding/json"
    "log"
    "net/http"
    "io/ioutil"
    "fmt"
//    "os"
//    "os/exec"
)

type test_struct struct {
    Test string
}

func packer(rw http.ResponseWriter, req *http.Request) {
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
    fmt.Fprintf(rw,"OK\n")
}

func main() {
    http.HandleFunc("/packer", packer)
    log.Fatal(http.ListenAndServe(":8080", nil))
}