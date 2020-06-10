package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/mikstew/mars-image/pkg/lrucache"
)

const (
	imageURL = "https://mars.com/image?lat=%g&long=%g"
	logLevel = log.DebugLevel
	// logLevel = log.WarnLevel
)

var cache = lrucache.Cache{}
var cacheNodes = map[string]*lrucache.Node{}
var cacheHit int64 = 0
var cacheMiss int64 = 0
var cacheEviction int64 = 0
var cacheHitRespTime int64 = 0
var cacheMissRespTime int64 = 0
var cacheEvictRespTime int64 = 0

// GetImageURL - Stubbed response. High cost call used after cache miss.
func GetImageURL(lat, long float64) string {
	return fmt.Sprintf("https://mars.com/image?lat=%g&long=%g", lat, long)
}

// image - executed when /image endpoint is called. Responds with image URL
// if query parameters meet validation. Calculates the following metrics:
// cache hits, cache misses, cache evictions, and average execution time
// for all three.
func image(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	var executionTime time.Duration = 0
	lat, long, err := parseParameters(r.URL.String())
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, err.Error())
		log.Debug(err.Error())
		return
	}
	key := fmt.Sprintf(imageURL, *lat, *long)
	log.Debug(key)
	if node, ok := cacheNodes[key]; ok {
		cache.MoveToMru(node)
		executionTime = time.Since(startTime) / time.Microsecond
		cacheHitRespTime = (int64(cacheHitRespTime)*cacheHit + int64(executionTime)) / (cacheHit + 1)
		cacheHit++
		log.WithFields(log.Fields{
			"Cache Hit Count":        cacheHit,
			"Execution Time (ms)":    executionTime,
			"Average Resp Time (ms)": cacheHitRespTime}).Debug("Cache hit.")
	} else {
		newKey := GetImageURL(*lat, *long)
		newNode, deletedNode := cache.AddToCache(newKey)
		if deletedNode != nil {
			log.Info("Cache eviction.")
			delete(cacheNodes, *deletedNode)
			executionTime = time.Since(startTime) / time.Microsecond
			cacheEvictRespTime = (int64(cacheEvictRespTime)*cacheEviction + int64(executionTime)) / (cacheEviction + 1)
			cacheEviction++
			log.WithFields(log.Fields{
				"Cache Eviction Count":   cacheEviction,
				"Execution Time (ms)":    executionTime,
				"Average Resp Time (ms)": cacheEvictRespTime}).Debug("Cache eviction.")
		}
		cacheNodes[newKey] = newNode
		executionTime = time.Since(startTime) / time.Microsecond
		cacheMissRespTime = (int64(cacheMissRespTime)*cacheMiss + int64(executionTime)) / (cacheMiss + 1)
		cacheMiss++
		log.WithFields(log.Fields{
			"Cache Miss Count":       cacheMiss,
			"Execution Time (ms)":    executionTime,
			"Average Resp Time (ms)": cacheMissRespTime}).Debug("Cache miss.")
	}
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", key)
}

// parseParameters - Validate the existance and value of query parameters.
// lat must be -90 to 90. long must be -180 to 180.
func parseParameters(s string) (*float64, *float64, error) {
	u, err := url.Parse(s)
	if err != nil {
		log.Println("Error parsing URL")
		return nil, nil, fmt.Errorf("Error parsing URL: %s", err)
	}
	q := u.Query()
	var lat float64 = 0
	var long float64 = 0
	if len(q["lat"]) > 0 {
		lat, err = strconv.ParseFloat(q["lat"][0], 64)
		if err != nil {
			log.Println("Error parsing 'lat' value from query")
			return nil, nil, fmt.Errorf("Error parsing 'lat' value: %s", err)
		}
		if lat < -90 || lat > 90 {
			log.Println("")
			return nil, nil, fmt.Errorf("Error parsing 'lat' value: must be -90 to 90")
		}
	} else {
		return nil, nil, fmt.Errorf("Query parameter 'lat' is required")
	}

	if len(q["long"]) > 0 {
		long, err = strconv.ParseFloat(q["long"][0], 64)
		if err != nil {
			return nil, nil, fmt.Errorf("Error parsing 'long' value: %s", err)
		}
		if long < -180 || long > 180 {
			return nil, nil, fmt.Errorf("Error parsing 'lat' value: must be -180 to 180")
		}
	} else {
		return nil, nil, fmt.Errorf("Query parameter 'long' is required")
	}

	return &lat, &long, nil
}

// getMetrics - executed when /get-metrics endpoint is called to
// report all metrics
func getMetrics(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w,
		"Cache Hits: %d - Average Response Time: %dms\n"+
			"Cache Misses: %d - Average Response Time: %dms\n"+
			"Cache Evictions: %d - Average Response Time: %dms",
		cacheHit, cacheHitRespTime,
		cacheMiss, cacheMissRespTime,
		cacheEviction, cacheEvictRespTime)
	log.Debug("/get-metrics endpoint called")
}

// resetMetrics - executed when /reset-metrics endpoint is called to
// reset all metrics
func resetMetrics(w http.ResponseWriter, r *http.Request) {
	cacheMiss = 0
	cacheHit = 0
	cacheEviction = 0
	cacheHitRespTime = 0
	cacheMissRespTime = 0
	cacheEvictRespTime = 0
	fmt.Fprintf(w, "Cache metrics reset.")
	log.Debug("/reset-metrics endpoint called")
}

// init - configure log settings
func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(logLevel)
}

func main() {
	http.HandleFunc("/image", image)
	http.HandleFunc("/get-metrics", getMetrics)
	http.HandleFunc("/reset-metrics", resetMetrics)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
