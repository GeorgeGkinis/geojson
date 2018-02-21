package main

import (
	"encoding/gob"
	"fmt"
	fb "geojson/featureBatch"
	"net"
	"os"
	"strconv"
	"sync"
)

func main() {
	fmt.Println("Starting Service B.")
	url := os.Getenv("SERVICE_B_URL")
	if url == "" {
		//fmt.Println("SERVICE_B_URL evnironment variable not set.\n Using default: 127.0.0.1:8080")
		url = "127.0.0.1:8080"
	}

	ln, err := net.Listen("tcp", url)
	if err != nil {
		fmt.Println("Cannot listen to: ", url, "Error: ", err)
	}
	for {
		conn, err := ln.Accept() // this blocks until connection or error
		if err != nil {
			fmt.Println("Error while recieving connection: ", err)
			continue
		}
		go receiveBatch(conn) // a goroutine handles conn so that the loop can accept other connections
	}
	fmt.Println("Exiting Service B.")
}

func receiveBatch(conn net.Conn) {
	// Create decoder listening on connection
	dec := gob.NewDecoder(conn)

	batch := &fb.FeatureBatch{}
	err := dec.Decode(batch)
	if err != nil {
		fmt.Println("Error while decoding batch: ", err)
	}
	conn.Close()

	fmt.Printf("Received batch %v of length: %v\n", batch.BatchNumber, len(batch.Features))
	batch.CalculatePopulationDensity()

	var wg sync.WaitGroup
	wg.Add(1)

	// Get retries from environment
	r, err := strconv.Atoi(os.Getenv("RETRIES"))
	if err != nil {
		//fmt.Println("RETRIES evnironment variable not set.\n Using default: 5")
		r = 5
	}

	url := os.Getenv("SERVICE_C_URL")
	if url == "" {
		//fmt.Println("SERVICE_C_URL evnironment variable not set.\n Using default: 127.0.0.1:8090")
		url = "127.0.0.1:8090"
	}
	go batch.Send(url, r, &wg)
	wg.Wait()
}
