package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"sync"
	"time"
)

// Request structure for incoming JSON payload
type Request struct {
	ToSort [][]int `json:"to_sort"`
}

// Response structure for outgoing JSON response
type Response struct {
	SortedArrays [][]int `json:"sorted_arrays"`
	TimeNS       int64   `json:"time_ns"`
}

func main() {
	http.HandleFunc("/process-single", processSingle)
	http.HandleFunc("/process-concurrent", processConcurrent)

	// log.Fatal(http.ListenAndServe(":8000", nil))
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func processSingle(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, sequentialSort)
}

func processConcurrent(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, concurrentSort)
}

func handleRequest(w http.ResponseWriter, r *http.Request, sorter func([][]int) ([][]int, int64)) {
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sortedArrays, timeNS := sorter(req.ToSort)

	resp := Response{
		SortedArrays: sortedArrays,
		TimeNS:       timeNS,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func sequentialSort(arrays [][]int) ([][]int, int64) {
	start := time.Now()
	result := make([][]int, len(arrays))
	for i, slice := range arrays {
		result[i] = append([]int(nil), slice...)
		sort.Ints(result[i])
	}
	return result, time.Since(start).Nanoseconds()
}

func concurrentSort(arrays [][]int) ([][]int, int64) {
	start := time.Now()
	result := make([][]int, len(arrays))
	wg := sync.WaitGroup{}

	for i, slice := range arrays {
		wg.Add(1)
		go func(idx int, s []int) {
			defer wg.Done()
			result[idx] = append([]int(nil), s...)
			sort.Ints(result[idx])
		}(i, slice)
	}

	wg.Wait()

	return result, time.Since(start).Nanoseconds()
}
