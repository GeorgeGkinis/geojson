package main

import (
	"encoding/gob"
	"fmt"
	fb "geojson/featureBatch"
	"net"
	"sync"
)

func main() {
	fmt.Println("Starting Service B.")
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		// handle error
	}
	for {
		conn, err := ln.Accept() // this blocks until connection or error
		if err != nil {
			// handle error
			fmt.Println("Error while recieving connection: ", err)
			continue
		}
		go handleConnection(conn) // a goroutine handles conn so that the loop can accept other connections
	}
	fmt.Println("Exiting Service C.")
}

func handleConnection(conn net.Conn) {
	// Create decoder listening on connection
	dec := gob.NewDecoder(conn)

	batch := &fb.FeatureBatch{}
	dec.Decode(batch)
	conn.Close()

	fmt.Printf("Received batch %v of length: %v\n", batch.BatchNumber, len(batch.Features))
	batch.CalculatePopulationDensity()

	var wg sync.WaitGroup
	wg.Add(1)
	go batch.Send("localhost:8090", 5, &wg)
	wg.Wait()
}
