// worker.go
package main

import (
	"log"
	"os/exec"
)

func worker(id int, works <-chan WorkRequest, results chan<- WorkEntry) {
	for j := range works {
		log.Printf("worker %d, processing job %s\n", id, j.UID)
		status := "process"
		RedisSet(WorkEntry{UID: j.UID, Status: status})
		cmd := exec.Command("/tmp/long", j.UID, "2")
		err := cmd.Start()
		if err != nil {
			status = "fail"
		}
		log.Println("Waiting for command to finish...")
		err = cmd.Wait()
		if err != nil {
			status = "fail"
		} else {
			status = "complete"
		}
		log.Printf("worker %d, job %s complete with status %s\n", id, j.UID, status)
		resp := WorkEntry{UID: j.UID, Status: status}
		RedisSet(resp)
		results <- resp
	}
}
