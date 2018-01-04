/*******************************************************************************
Copyright (c) 2018 Venish Joe Clarence (http://venishjoe.net)

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

******************************************************************************/

package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/tidwall/gjson"
)

//LocationHeatMapData - Data structure to hold heat map data
type LocationHeatMapData struct {
	Latitude, Longitude string
	Weight              float64
}

//LocationHeatMap - Data structure to hold LocationHeatMapData
type LocationHeatMap struct {
	LocationHeatMap []LocationHeatMapData
}

func main() {

	//Read commandline arguments
	disableTLSSecurity := flag.Bool("disableTLSSecurity", false, "To disable TLS Security, send true")
	flag.Parse()

	//Disable TLS Security if disableTLSSecurity is true
	httpClient := &http.Client{}
	if *disableTLSSecurity {
		tlsWithSecDisabled := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		httpClient = &http.Client{Transport: tlsWithSecDisabled}
	}

	//Read Google API Key
	googleAPIKey, errReadingAPIKey := ioutil.ReadFile("../data/googlegeoapi.key")
	handleFatalError(errReadingAPIKey)

	//Build Web Service URL
	googleAPIBaseURL := "https://maps.googleapis.com/maps/api/geocode/json?key=" + string(googleAPIKey) + "&address="

	//Read Input File with Locations Data
	fileWithLocationData, errReadingLocationData := os.Open("../data/location.data")
	handleFatalError(errReadingLocationData)
	defer fileWithLocationData.Close()
	fileWithLocationDataScanner := bufio.NewScanner(fileWithLocationData)

	//Create file to write output delimited data
	locationDataFileWithLatLng, errWritingFileWithLatLng := os.Create("../data/locationlatlng.data")
	handleFatalError(errWritingFileWithLatLng)
	defer locationDataFileWithLatLng.Close()

	//Create file to write output JS for Goolgle Maps Rendering
	jsDataForHeatMap, errWritingJSDataForHeatMap := os.Create("../ui/places_data.js")
	handleFatalError(errWritingJSDataForHeatMap)
	defer jsDataForHeatMap.Close()

	//Create Heat Map Data Java Script (places_data.js)
	fmt.Fprintf(jsDataForHeatMap, "function getDataPoints() { \n\treturn [\n")
	var locationHeatMap LocationHeatMap

	//Iterate location data file and process records
	log.Println("Processing locations........")
	for fileWithLocationDataScanner.Scan() {
		googleAPIConstructedURL := googleAPIBaseURL + url.QueryEscape(fileWithLocationDataScanner.Text())

		response, errFromAPI := httpClient.Get(googleAPIConstructedURL)
		handleFatalError(errFromAPI)
		defer response.Body.Close()

		//Process response from service if HTTP status is OK
		if response.StatusCode == http.StatusOK {
			responseBytes, errFromReadResponse := ioutil.ReadAll(response.Body)
			handleFatalError(errFromReadResponse)

			responseString := string(responseBytes)
			jsonStatusCode := gjson.Get(responseString, "status")

			//Process response from service if JSON status is OK
			if jsonStatusCode.String() == "OK" {
				indexInner := 0
				formattedAddressArray := gjson.Get(responseString, "results.#.formatted_address")
				for _, formattedAddress := range formattedAddressArray.Array() {
					fmt.Fprintf(locationDataFileWithLatLng, formattedAddress.String()+"$")
					addressComponentsArray := gjson.Get(responseString, "results."+strconv.Itoa(indexInner)+".address_components.#.short_name")

					for _, addressComponentsArray := range addressComponentsArray.Array() {
						fmt.Fprintf(locationDataFileWithLatLng, addressComponentsArray.String()+"$")
					}
					geometryLocationLat := gjson.Get(responseString, "results."+strconv.Itoa(indexInner)+".geometry.location.lat")
					geometryLocationLng := gjson.Get(responseString, "results."+strconv.Itoa(indexInner)+".geometry.location.lng")
					fmt.Fprintf(locationDataFileWithLatLng, geometryLocationLat.String()+"$"+geometryLocationLng.String()+"\n")

					if len(locationHeatMap.LocationHeatMap) == 0 {
						locationHeatMap.LocationHeatMap = append(locationHeatMap.LocationHeatMap,
							LocationHeatMapData{Latitude: geometryLocationLat.String(), Longitude: geometryLocationLng.String(), Weight: 0.5})
					} else {
						dataMatchIndex := checkIfValueExists(locationHeatMap, geometryLocationLat.String(), geometryLocationLng.String())
						if dataMatchIndex != -1 {
							locationHeatMap.LocationHeatMap[dataMatchIndex].Weight = locationHeatMap.LocationHeatMap[dataMatchIndex].Weight + 0.1
						} else {
							locationHeatMap.LocationHeatMap = append(locationHeatMap.LocationHeatMap,
								LocationHeatMapData{Latitude: geometryLocationLat.String(), Longitude: geometryLocationLng.String(), Weight: 0.5})
						}
					}
					indexInner++
				}
			} else {
				log.Println("JSON Error!")
			}
		}
		googleAPIConstructedURL = ""
	}
	log.Println("Completed")
	log.Println("Generating places JS........")
	//Update Heat Map Data Java Script (places_data.js)
	for indexLocationDataWriteIndex := 0; indexLocationDataWriteIndex < len(locationHeatMap.LocationHeatMap); indexLocationDataWriteIndex++ {
		fmt.Fprintf(jsDataForHeatMap, "\t\t{location: new google.maps.LatLng("+
			locationHeatMap.LocationHeatMap[indexLocationDataWriteIndex].Latitude+", "+
			locationHeatMap.LocationHeatMap[indexLocationDataWriteIndex].Longitude+"), weight: "+
			strconv.FormatFloat(locationHeatMap.LocationHeatMap[indexLocationDataWriteIndex].Weight, 'f', 1, 64)+"}")

		if indexLocationDataWriteIndex == len(locationHeatMap.LocationHeatMap)-1 {
			fmt.Fprintf(jsDataForHeatMap, "\n")
		} else {
			fmt.Fprintf(jsDataForHeatMap, ",\n")
		}
	}
	fmt.Fprintf(jsDataForHeatMap, "\t]; \n}")
	log.Println("Completed")
}

//Function to handle fatal errors
func handleFatalError(fatalError error) {
	if fatalError != nil {
		log.Fatal(fatalError)
		return
	}
}

//Function to check if Latitue/Longitue exists in data structure
func checkIfValueExists(locationHeatMapCheck LocationHeatMap, geometryLocationLatCheck, geometryLocationLngCheck string) int {
	for locationDataIndexCheck := 0; locationDataIndexCheck < len(locationHeatMapCheck.LocationHeatMap); locationDataIndexCheck++ {
		if locationHeatMapCheck.LocationHeatMap[locationDataIndexCheck].Latitude == geometryLocationLatCheck &&
			locationHeatMapCheck.LocationHeatMap[locationDataIndexCheck].Longitude == geometryLocationLngCheck {
			return locationDataIndexCheck
		}
	}
	return -1
}
