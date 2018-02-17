package main

import (
	geojson "../src/github.com/paulmach/go.geojson"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
)

// declare structure for country
type Country struct {
	Name              string
	Area              int
	Population        int
	PopulationDensity float32
}

func (c Country) String() string {

	// TODO: Implement some proper padding to align the fields in order to make output more presentable.
	return fmt.Sprintf("Population Density : %v,\t\t\t\t Area : %v,\t\t Population : %v,\t\t %v", c.PopulationDensity, c.Area, c.Population, c.Name)
}

func main() {

	//respBytes,err := getGeoJSON("https://ccdnn.locsensads.com/jobs/worldborders.geojson")
	respBytes, err := getGeoJSONFile("countries.geojson")
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	fc, err := geojson.UnmarshalFeatureCollection(respBytes)
	if err != nil {
		fmt.Println(err)
		os.Exit(-3)
	}

	var countries []Country

	for _, f := range fc.Features {

		name, err := f.PropertyString("NAME")
		if err != nil {
			println(err)
		}

		// TODO: If area == 0 then calculate from MultiPolygon.
		/*		https://gis.stackexchange.com/questions/711/how-can-i-measure-area-from-geographic-coordinates
				var area = 0.0;
				var len = ring.components && ring.components.length;
				if (len > 2) {
					var p1, p2;
					for (var i=0; i<len-1; i++) {
					p1 = ring.components[i];
					p2 = ring.components[i+1];
					area += OpenLayers.Util.rad(p2.x - p1.x) *
					(2 + Math.sin(OpenLayers.Util.rad(p1.y)) +
					Math.sin(OpenLayers.Util.rad(p2.y)));
					}
					area = area * 6378137.0 * 6378137.0 / 2.0;
				}*/

		area, err := f.PropertyInt("AREA")
		if err != nil {
			println(err)
		}

		pop, err := f.PropertyInt("POP2005")
		if err != nil {
			println(err)
		}

		countries = append(countries, Country{name, area, pop, 0.0})
	}

	sendToB(countries)

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

func sendToB(countries []Country) (err error) {
	err = nil

	serviceB(countries)

	return err
}

func serviceB(countries []Country) {

	for i, c := range countries {
		if c.Area != 0 {
			c.PopulationDensity = float32(c.Population) / float32(c.Area)
		} else {
			c.PopulationDensity = -1
		}
		countries[i] = c
	}

	sendToC(countries)
}

func sendToC(countries []Country) (err error) {
	err = nil
	serviceC(countries)
	return err
}

func serviceC(countries []Country) {
	sort.Slice(countries, func(i, j int) bool {
		return countries[i].PopulationDensity > countries[j].PopulationDensity
	})

	for _, c := range countries {
		fmt.Println(c)
	}
}
