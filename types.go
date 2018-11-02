package main

import "net/http"

type CityTemp struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	ID   int    `json:"id"`
	Name string `json:"name"`
	Cod  int    `json:"cod"`
}

type Config struct {
	city     string
	unit     string
	credFile string
	appID    string
	client   http.Client
}
