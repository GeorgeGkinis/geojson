package main

import (
	"encoding/gob"
	"fmt"
	geojson "github.com/paulmach/go.geojson"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

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

	// Send Features in batches.
	var batch []geojson.Feature
	batchSize := 5

	// Include total number of countries so that Service C knows when done receiving.
	header := new(geojson.Feature)
	header.SetProperty("NUM_OF_COUNTRIES", len(fc.Features))
	batch = append(batch, *header)

	// Fill batches and send to Service B
	for i, f := range fc.Features {

		batch = append(batch, *f)
		if len(batch) == batchSize || len(fc.Features) == i+1 {

			fmt.Printf("Sending batch: %+v, number of elements: %+v\n", 1+(i/batchSize), len(batch))
			sendToB(batch)
			batch = nil
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

func sendToB(batch []geojson.Feature) {
	// Connect to Service B
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("Connection error", err)
	}
	defer conn.Close()

	// Create encoder to send Features over the wire
	encoder := gob.NewEncoder(conn)

	// Send batch
	encoder.Encode(batch)
}
