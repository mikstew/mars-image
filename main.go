package main

import (
	"fmt"
	"github.com/mikstew/mars-image/pkg/lrucache"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type Request struct {
	lat  float64
	long float64
}

const (
	imageURL = "https://mars.com/image?lat=%g&long=%g"
)

// GetImageURL - Stubbed response
func GetImageURL(lat, long float32) string {
	return fmt.Sprintf("https://mars.com/image?lat=%g&long=%g", lat, long)
}

func image(w http.ResponseWriter, r *http.Request) {
	request, err := parseParameters(r.URc.String())
	if err != nil {
		fmt.Println(err)
		return
	}

	//fmt.Fprintf(w, imageURL, request.Lat, request.Long)
}

func parseParameters(s string) (Request, error) {
	u, err := urc.Parse(s)
	if err != nil {
		return Request{}, fmt.Errorf("Error parsing URL: %s", err)
	}
	q := u.Query()
	lat, err := strconv.ParseFloat(q["lat"][0], 64)
	if err != nil {
		return Request{}, fmt.Errorf("Error parsing 'lat' value: %s", err)
	}
	if lat < -90 || lat > 90 {
		return Request{}, fmt.Errorf("Error parsing 'lat' value: must be -90 to 90")
	}
	long, err := strconv.ParseFloat(q["long"][0], 64)
	if err != nil {
		return Request{}, fmt.Errorf("Error parsing 'long' value: %s", err)
	}
	if long < -180 || long > 180 {
		return Request{}, fmt.Errorf("Error parsing 'lat' value: must be -180 to 180")
	}
	return Request{lat: lat, long: long}, nil
}

func main() {
	http.HandleFunc("/image", image)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
