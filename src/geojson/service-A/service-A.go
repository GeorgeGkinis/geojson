package main

import (
	"encoding/gob"
	"fmt"
	geojson "github.com/paulmach/go.geojson"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

type featureBatch struct {
	Timestamp     int64
	TotalMessages int
	Features      []geojson.Feature
}

func main() {

	// Get GeoJSON
	//respBytes,err := getGeoJSON("https://ccdnn.locsensads.com/jobs/worldborders.geojson")
	respBytes, err := getGeoJSONFile("countriesSample.geojson")
	if err != nil {
		log.Fatal("Error reading geojson: ", err)
	}

	// Make an FeatureCollection out of it
	fc, err := geojson.UnmarshalFeatureCollection(respBytes)
	if err != nil {
		log.Fatal("Error while unmarshaling geojson: ", err)
	}

	batchSize := 5

	fb := featureBatch{
		Timestamp:     time.Now().UnixNano(),
		TotalMessages: len(fc.Features),
		Features:      nil,
	}
	// Fill batches and send to Service B
	for i, f := range fc.Features {

		fb.Features = append(fb.Features, *f)
		if len(fb.Features) == batchSize || len(fc.Features) == i+1 {

			fmt.Printf("Sending batch: %+v, number of elements: %+v\n", 1+(i/batchSize), len(fb.Features))
			sendToB(fb)
			fb.Features = nil
		}
	}
	fmt.Println("done")
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

	return ioutil.ReadFile(path)
}

func sendToB(fb featureBatch) {
	// Connect to Service B
	conn, err := net.Dial("tcp", "localhost:8080")
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
	//fmt.Println(fb)
}
