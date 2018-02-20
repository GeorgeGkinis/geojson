package main

import (
	"fmt"
	fb "geojson/featureBatch"
	geojson "github.com/paulmach/go.geojson"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {

	fmt.Println("Starting Service A.")

	// Set batch size
	batchSize := 5

	// Get GeoJSON
	//respBytes,err := getGeoJSON("https://ccdnn.locsensads.com/jobs/worldborders.geojson")
	respBytes, err := getGeoJSONFile("countries.geojson")
	if err != nil {
		log.Fatal("Error reading geojson: ", err)
	}

	// Make an FeatureCollection out of GeoJSON
	fc, err := geojson.UnmarshalFeatureCollection(respBytes)
	if err != nil {
		log.Fatal("Error while unmarshaling geojson: ", err)
	}

	batch := fb.FeatureBatch{
		Timestamp:     time.Now().UnixNano(),
		BatchNumber:   0,
		TotalMessages: len(fc.Features),
		Features:      nil,
	}
	// Fill batches and send to Service B
	var wg sync.WaitGroup
	for i, f := range fc.Features {

		batch.Features = append(batch.Features, *f)

		// If batch is full OR no remaining Features to send
		if len(batch.Features) == batchSize || len(fc.Features) == i+1 {
			// Set batch number
			batch.BatchNumber = 1 + (i / batchSize)

			wg.Add(1)
			go batch.Send("localhost:8080", 5, &wg)

			// clear batch.Features for reuse
			batch.Features = nil
		}
	}
	// Wait for all batches to finish sending.
	wg.Wait()
	fmt.Println("Exiting Service A.")
}

func getGeoJSON(url string) ([]byte, error) {

	// get geojson file
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Could not retrieve GeoJSON from " + url + ". Status code is : " + string(resp.Status))
		//os.Exit(-1)
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		//os.Exit(-2)
	}

	return respBytes, err
}

func getGeoJSONFile(path string) ([]byte, error) {
	b, err := ioutil.ReadFile(path)
	return b, err
}
