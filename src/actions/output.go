package actions

import (
	"encoding/json"
	. "github.com/lcserny/goutils"
)

type OutputRequestData struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	MoviesPath string `json:"moviesPath"`
	TvPath     string `json:"tvPath"`
}

type OutputResponseData struct {
	Name string `json:"name"`
}

// use tmdbApiKey global var
func OutputAction(jsonPayload []byte) (string, error) {
	// TODO: based on data provided
	//  trim name and such, make it decent and find it on disk
	//  if nothing matched on disk, get a list of 10? API resolved output names from that decent name
	//  if nothing returned from API, return the decent name (listOf 1)

	source := "Fantstic BeastThe Crimes of Grindwald"
	target := "Fantastic Beasts: The Crimes of Grindelwald"

	distance := LevenshteinDistance(source, target)

	resultBytes, _ := json.Marshal(distance)

	return string(resultBytes), nil
}
