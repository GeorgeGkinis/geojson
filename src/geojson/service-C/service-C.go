package main

import (
	"encoding/gob"
	"fmt"
	geojson "github.com/paulmach/go.geojson"
	"log"
	"net"
	"sort"
	_ "sort"
)

type featureBatch struct {
	Timestamp     int64
	TotalMessages int
	Features      []geojson.Feature
}

func main() {
	fmt.Println("Starting Service C")

	queue := make(chan featureBatch)

	go listen(queue)

	var features []geojson.Feature

	// Read incomming featureBatches
	for fb := range queue {
		for _, f := range fb.Features {
			features = append(features, f)
		}

		// When all features are received create an FeatureCollection
		if len(features) == fb.TotalMessages {
			// Create closure to keep a copy of features
			func([]geojson.Feature) {
				fc := geojson.NewFeatureCollection()

				for _, f := range features {
					a := f
					fc.AddFeature(&a)
				}

				// Sort countries based on POPDENS
				sortFeatures(fc)
				// Print countries on screen
				//printfeatures(fc)
			}(features)
		}
	}
}

func listen(queue chan featureBatch) {
	ln, err := net.Listen("tcp", ":8090")
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
		go handleConnection(conn, queue) // a goroutine handles conn so that the loop can accept other connections
	}
}

func handleConnection(conn net.Conn, queue chan featureBatch) {
	// Create decoder listening on connection
	dec := gob.NewDecoder(conn)

	fb := &featureBatch{}
	if err := dec.Decode(fb); err != nil {
		fmt.Println("Something went wrong while receiving batch: ", err)
	}
	conn.Close()

	//fmt.Printf("Received : %+v\n", fb)

	queue <- *fb
}

func sortFeatures(fc *geojson.FeatureCollection) {
	sort.Slice(fc.Features, func(i, j int) bool {

		idens, err := fc.Features[i].PropertyFloat64("POPDENS")
		if err != nil {
			log.Fatal("Error while reading POPDENS: ", err)
		}

		jdens, err := fc.Features[j].PropertyFloat64("POPDENS")
		if err != nil {
			log.Fatal("Error while reading POPDENS: ", err)
		}
		return idens > jdens
	})
	printfeatures(fc)
}

func printfeatures(fc *geojson.FeatureCollection) {
	for _, f := range fc.Features {
		d, _ := f.PropertyFloat64("POPDENS")
		n, _ := f.PropertyString("NAME")
		fmt.Printf("%9.f\t\t%+v\n", d, n)
	}
}
