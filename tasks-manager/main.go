package main

import (
	"sadath/client"
	"sadath/scheduler"
	"sadath/server"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(3)
	go client.Client(&wg)
	go server.Workers(&wg)
	go scheduler.Scheduler(&wg)
	wg.Wait()
}
