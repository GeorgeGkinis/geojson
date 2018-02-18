package main

import (
	"encoding/gob"
	"fmt"
	geojson "github.com/paulmach/go.geojson"
	"log"
	"net"
)

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

	batch := &[]geojson.Feature{}
	dec.Decode(batch)
	conn.Close()

	fmt.Println("Received batch of length: ", len(*batch))
	calculatePopulationDensity(*batch)

	for _, f := range *batch {
		n, _ := f.PropertyString("NAME")
		fmt.Printf("Sent: %+v\n", n)
	}
	sendToC(*batch)
}

func calculatePopulationDensity(f []geojson.Feature) {

	for _, c := range f {

		// Do not process header feature.
		if _, err := c.PropertyInt("NUM_OF_COUNTRIES"); err != nil {

			// Get polulation
			pop, err := c.PropertyInt("POP2005")
			if err != nil {
				log.Fatal(err)
			}

			// Get Area
			area, err := c.PropertyInt("AREA")
			if err != nil {
				log.Fatal(err)
			}

			if area != 0 {
				dens := float64(pop) / float64(area)
				c.SetProperty("POPDENS", dens)
			} else {
				// TODO: If AREA == 0 then calculate from multipolygon
				c.SetProperty("POPDENS", -1)
			}
		}
	}
}

func sendToC(batch []geojson.Feature) {
	// Connect to Service C
	conn, err := net.Dial("tcp", "localhost:8090")
	if err != nil {
		log.Fatal("Connection error", err)
	}
	defer conn.Close()

	// Create encoder to send Features over the wire
	encoder := gob.NewEncoder(conn)

	// Send batch
	encoder.Encode(batch)
}
