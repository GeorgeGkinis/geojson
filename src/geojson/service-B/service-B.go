package main

import (
	"encoding/gob"
	"fmt"
	geojson "github.com/paulmach/go.geojson"
	"log"
	"net"
)

type featureBatch struct {
	Timestamp     int64
	TotalMessages int
	Features      []geojson.Feature
}

func main() {
	fmt.Println("Starting Service B")
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
}

func handleConnection(conn net.Conn) {
	// Create decoder listening on connection
	dec := gob.NewDecoder(conn)

	fb := &featureBatch{}
	dec.Decode(fb)
	conn.Close()

	fmt.Println("Received batch of length: ", len(fb.Features))
	calculatePopulationDensity(fb)

	sendToC(fb)
}

func calculatePopulationDensity(fb *featureBatch) {

	for _, f := range fb.Features {

		// Get polulation
		pop, err := f.PropertyInt("POP2005")
		if err != nil {
			log.Fatal(err)
		}

		// Get Area
		area, err := f.PropertyInt("AREA")
		if err != nil {
			log.Fatal(err)
		}

		if area != 0 {
			dens := float64(pop) / float64(area)
			f.SetProperty("POPDENS", dens)
		} else {
			// TODO: If AREA == 0 then calculate from multipolygon
			f.SetProperty("POPDENS", -1)
		}
	}

}

func sendToC(fb *featureBatch) {
	// Connect to Service C
	conn, err := net.Dial("tcp", "localhost:8090")
	if err != nil {
		log.Fatal("Connection error", err)
	}
	defer conn.Close()

	// Create encoder to send Features over the wire
	encoder := gob.NewEncoder(conn)

	// Send batch
	if err = encoder.Encode(fb); err != nil {
		fmt.Println("Something went wrong while sending batch: ", err)
	}
}
