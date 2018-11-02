package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const SRV = "https://api.openweathermap.org/data/2.5/weather"

func main() {

	var c Config

	flag.StringVar(&c.city, "c", "Leipzig,DE", "Query whether for City")
	flag.StringVar(&c.credFile, "f", "./APPID", "File containing APPID key")
	flag.StringVar(&c.unit, "u", "celsius", "Units for temperature (celsius, fahrenheit)")
	flag.Parse()

	var metric string
	switch c.unit {
	case "celsius":
		metric = "metric"
	case "fahrenheit":
		metric = "imperial"
	default:
		log.Fatalf("unsupported unit type: %s", c.unit)
	}

	// get appID key to authenticate against OpenWeather API
	f, err := os.Open(c.credFile)
	if err != nil {
		log.Fatalf("could not open credentials file: %v", err)
	}

	d, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("could not read credentials: %v", err)
	}

	c.appID = string(d)

	// generate request
	req, err := getRequest(c, metric)
	if err != nil {
		log.Fatalf("could not create request: %v", err)
	}

	// create client
	c.client = newClient()

	// send request
	resp, err := c.client.Do(req)
	if err != nil {
		log.Fatalf("could not get response: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("HTTP status code error: %s", resp.Status)
	}

	var t = CityTemp{}
	err = json.NewDecoder(resp.Body).Decode(&t)
	if err != nil {
		log.Fatalf("could not decode data: %v", err)
	}

	fmt.Printf("Temperature in %s: %0.2f %s", t.Name, t.Main.Temp, c.unit)
	fmt.Println()

}

func getRequest(c Config, metric string) (*http.Request, error) {
	req, err := http.NewRequest("GET", SRV, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("q", c.city)
	q.Add("units", metric)
	q.Add("APPID", c.appID)
	req.URL.RawQuery = q.Encode()

	return req, nil
}

func newClient() http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := http.Client{
		Timeout:   3 * time.Second,
		Transport: tr,
	}

	return client
}
