package featureBatch

import (
	"encoding/gob"
	"fmt"
	"github.com/paulmach/go.geojson"
	"log"
	"net"
	"sync"
	"time"
)

type FeatureBatch struct {
	Timestamp     int64
	TotalMessages int
	BatchNumber   int
	Features      []geojson.Feature
}

func (b FeatureBatch) Send(url string, attempts int, wg *sync.WaitGroup) {
	fmt.Printf("Sending batch: %+v, number of elements: %+v\n", b.BatchNumber, len(b.Features))
	defer wg.Done()
	for i := 1; ; i++ {
		//var err *error
		err := b.sendBatch(url)

		if err == nil {
			return
		}

		if i >= (attempts) {
			fmt.Printf("Failed batch %v after %d attempts, last error: %v\n", b.BatchNumber, i, err)
			break
		}

		fmt.Printf("Something went wrong while sending batch %v: %v\n", b.BatchNumber, err)
		fmt.Printf("Retrying in %+v  second(s)..\n", i)
		time.Sleep(time.Duration(i) * time.Second)
	}
}

func (b FeatureBatch) sendBatch(url string) error {

	// Connect to Service B
	conn, err := net.Dial("tcp", url)

	if err != nil {
		//fmt.Println("Connection error: ", err)
		return err
	}
	defer conn.Close()

	// Create encoder to send Features over the wire
	encoder := gob.NewEncoder(conn)

	err = encoder.Encode(b)
	return err
}

func (b FeatureBatch) CalculatePopulationDensity() {

	for _, f := range b.Features {

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
			f.SetProperty("POPDENS", 0.0)
		}
	}

}
