package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/golang/gddo/httputil/header"
	"github.com/gorilla/mux"
)

type flight struct {
	Source string `json:"source"`
	Dest   string `json:"dest"`
}

func main() {
	//flightrouter := mux.NewRouter().StrictSlash(true)
	flightrouter := mux.NewRouter()
	flightrouter.HandleFunc("/calculate", findFlightPath).Methods("POST")
	log.Fatal(http.ListenAndServe(":8081", flightrouter))
}

func findFlightPath(w http.ResponseWriter, r *http.Request) {
	flightData := []flight{}

	// decode request from HTTP to check if the input is JSON or not.
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			http.Error(w, "Not JSON Header", http.StatusUnsupportedMediaType)
			return
		}
	}
	// read data from the json file max - 1MB
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	req, _ := ioutil.ReadAll(r.Body)

	// Unmarshall Data now into src,dest pairs (struct f)
	// If the data is not formatted in src, dest pairs, return error
	if err := json.Unmarshal(req, &flightData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// If Input data is null, return BadRequest
	if len(flightData) == 0 {
		http.Error(w, "Error: NULL Data", http.StatusBadRequest)
		return
	}
	//cityarr := [][]string{{"IND", "EWR"}, {"SFO", "ATL"}, {"GSO", "IND"}, {"ATL", "GSO"}}
	//cityarr := [][]string{{"ATL", "EWR"}, {"SFO", "ATL"}}
	// We have the data in above form, do a topology Sort and then return the output
	res, err := topologySort(flightData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Encode output and send it back
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

/*
 * Add a Graph to find all city pairs and the inDegree of each city.
 * if more than one origin city, return error.
 * Otherwise, find the start and destination cities and return that data.
 */
func topologySort(citypairs []flight) (flight, error) {
	citymap := map[string][]string{}
	inDegree := map[string]int{}
	cn := len(citypairs)
	for i := 0; i < cn; i++ {
		v := citypairs[i]
		citymap[v.Source] = append(citymap[v.Source], v.Dest)
		inDegree[v.Dest] += 1
	}

	respath := []string{}
	stk := []string{}
	res := flight{}
	count := 0
	for k := range citymap {
		v, _ := inDegree[k]
		if v == 0 {
			stk = append(stk, k)
			count++
		}
	}

	// check for data input errors; multiple origin cities
	if count > 1 {
		return res, errors.New("ERROR: Multiple Origin Cities")
	}
	for len(stk) > 0 {
		city := stk[0]
		respath = append(respath, city)
		stk = stk[1:]
		v, _ := citymap[city]
		for i := 0; i < len(v); i++ {
			inDegree[v[i]] -= 1
			kv, _ := inDegree[v[i]]
			if kv == 0 {
				stk = append(stk, v[i])
			}
		}
	}
	// Other error handling...
	// Unvisited cities and cycles
	for _, v := range inDegree {
		if v > 0 {
			return res, errors.New("ERROR: Unvisited Cities, Input Error")
		}
	}

	res.Source = respath[0]
	res.Dest = respath[len(respath)-1]
	return res, nil
}
