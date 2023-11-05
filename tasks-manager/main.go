package main

import (
	"sadath/client"
	"sadath/server"
	  "sync"

)

func main(){
	var wg sync.WaitGroup
	wg.Add(2)
	go client.Client(&wg)
	go server.Workers(&wg)
	wg.Wait()
}