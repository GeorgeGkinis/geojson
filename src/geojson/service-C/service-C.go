package main

import (
	"encoding/gob"
	"fmt"
	geojson "github.com/paulmach/go.geojson"
	"net"
	_ "sort"
)

func main() {
	fmt.Println("Starting Service C")
	ln, err := net.Listen("tcp", ":8090")
	if err != nil {
		// handle error
	}

	fc := new([]geojson.Feature)
	fcSize := -1

	for {
		conn, err := ln.Accept() // this blocks until connection or error
		if err != nil {
			// handle error
			fmt.Println("Error while recieving connection: ", err)
			continue
		}
		go handleConnection(conn, fc, &fcSize) // a goroutine handles conn so that the loop can accept other connections
	}

}

func handleConnection(conn net.Conn, fc *[]geojson.Feature, fcSize *int) {
	// Create decoder listening on connection
	dec := gob.NewDecoder(conn)

	batch := &[]geojson.Feature{}
	dec.Decode(batch)
	conn.Close()

	fmt.Printf("Received : %+v\n", batch)

	for _, f := range *batch {

		// Register total number of countries to recieve
		if n, err := f.PropertyInt("NUM_OF_COUNTRIES"); err == nil {
			*fcSize = n
		} else {
			*fc = append(*fc, f)

		}
	}

	/*	if len(*fc) == *fcSize {

		// sort countries by population density
		sort.Slice(fc, func(i, j int) bool {

			idens,err := fc[i].PropertyFloat64("POPDENS")
			if err != nil {log.Fatal("Error while reading POPDENS: ",err)}

			kdens,err := fc.Features[i].PropertyFloat64("POPDENS")
			if err != nil {log.Fatal("Error while reading POPDENS: ",err)}

			return idens>kdens
		})*/

	for _, f := range *fc {

		d, _ := f.PropertyFloat64("POPDENS")
		n, _ := f.PropertyString("NAME")
		fmt.Printf("%+v\t\t%+v\n", d, n)
	}
	// TODO: draw a map
}
