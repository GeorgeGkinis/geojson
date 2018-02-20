package main

import (
	"encoding/gob"
	"fmt"
	fb "geojson/featureBatch"
	geojson "github.com/paulmach/go.geojson"
	"github.com/rs/cors"
	"log"
	"net"
	"net/http"
	"sort"
)

func main() {
	fmt.Println("Starting Service C.")

	queue := make(chan fb.FeatureBatch)

	// Slice to hold batches
	var features []geojson.Feature

	// Serve at a different goroutine
	go listen(queue)

	// Create FeatureCollection to serve at endpoint.
	fc := geojson.NewFeatureCollection()

	go serve(fc)

	// Read incomming featureBatches from queue
	for batch := range queue {
		for _, f := range batch.Features {
			features = append(features, f)
		}

		// When all features are received
		if len(features) == batch.TotalMessages {

			// create a FeatureCollection
			tmpfc := geojson.NewFeatureCollection()

			for _, f := range features {
				a := f
				tmpfc.AddFeature(&a)
			}
			// Sort countries based on POPDENS
			sortByDensity(tmpfc)

			// Assign completed FeatureCollection to fc.
			// This way only completed FeatureCollections are served.
			fc.Features = tmpfc.Features

			// Print countries on screen
			printfeatures(fc)

			// Clear features slice for reuse
			features = nil
		}
	}
	close(queue)
	fmt.Println("Exiting Service C.")
}

func listen(queue chan fb.FeatureBatch) {
	ln, err := net.Listen("tcp", "localhost:8090")
	if err != nil {
		fmt.Println("Error setting up connection: " + err.Error())
	}
	for {
		conn, err := ln.Accept() // this blocks until connection or error
		if err != nil {
			// handle error
			fmt.Println("Error while recieving connection: ", err)
			continue
		}
		handleConnection(conn, queue) // a goroutine handles conn so that the loop can accept other connections
	}
}

func serve(fc *geojson.FeatureCollection) {

	mux := http.NewServeMux()
	mux.HandleFunc("/countries", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Disposition", "attachment; filename=\"countries.geojson\"")

		b, err := fc.MarshalJSON()
		if err != nil {
			fmt.Println("Error marshaling GeoJSON: ", err)
		}
		w.Write(b)
	})

	mux.Handle("/", http.FileServer(http.Dir("src/geojson/service-C")))

	handler := cors.Default().Handler(mux)
	http.ListenAndServe(":8091", handler)
}

func handleConnection(conn net.Conn, queue chan fb.FeatureBatch) {
	// Create decoder listening on connection
	dec := gob.NewDecoder(conn)

	batch := &fb.FeatureBatch{}
	if err := dec.Decode(batch); err != nil {
		fmt.Println("Something went wrong while receiving batch: ", err)
	}
	conn.Close()

	//fmt.Printf("Received : %+v\n", fb)

	queue <- *batch
}

func sortByDensity(fc *geojson.FeatureCollection) {
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
}

func printfeatures(fc *geojson.FeatureCollection) {
	for _, f := range fc.Features {
		d, _ := f.PropertyFloat64("POPDENS")
		n, _ := f.PropertyString("NAME")
		fmt.Printf("%9.f\t\t%+v\n", d, n)
	}
}
