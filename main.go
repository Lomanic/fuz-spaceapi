package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Config struct {
	PORT        string
	PRESENCEAPI string
	SPACEAPI    string
}

type SpaceAPI struct {
	API              string   `json:"api"`
	APICompatibility []string `json:"api_compatibility,omitempty"`
	Space            string   `json:"space"`
	Logo             string   `json:"logo"`
	URL              string   `json:"url"`
	Location         struct {
		Address string  `json:"address,omitempty"`
		Lat     float64 `json:"lat"`
		Lon     float64 `json:"lon"`
	} `json:"location"`
	Contact struct {
		Email   string `json:"email,omitempty"`
		Irc     string `json:"irc,omitempty"`
		Ml      string `json:"ml,omitempty"`
		Twitter string `json:"twitter,omitempty"`
		Matrix  string `json:"matrix,omitempty"`
	} `json:"contact"`
	IssueReportChannels []string `json:"issue_report_channels"`
	State               struct {
		Icon struct {
			Open   string `json:"open"`
			Closed string `json:"closed"`
		} `json:"icon,omitempty"`
		Open       *bool   `json:"open"`
		Message    string `json:"message,omitempty"`
		LastChange int64  `json:"lastchange,omitempty"`
	} `json:"state"`
	Projects []string `json:"projects,omitempty"`
}

type Status struct {
	FuzIsOpen      bool      `json:"fuzIsOpen"`
	LastSeenAsOpen bool      `json:"lastSeenAsOpen"`
	LastSeen       time.Time `json:"lastSeen"`
	LastOpened     time.Time `json:"lastOpened"`
	LastClosed     time.Time `json:"lastClosed"`
}

var (
	config = Config{
		PORT: "8080",
	}
	spaceAPI = SpaceAPI{
		API: "0.13",
	}
)

func init() {
	port := os.Getenv("PORT")
	if val, _ := strconv.Atoi(port); val > 0 {
		config.PORT = port
	}
	config.PRESENCEAPI = os.Getenv("PRESENCEAPI")
	if config.PRESENCEAPI == "" {
		panic("PRESENCEAPI is empty")
	}
	config.SPACEAPI = os.Getenv("SPACEAPI")
	if config.SPACEAPI == "" {
		panic("SPACEAPI is empty")
	}
	err := json.Unmarshal([]byte(config.SPACEAPI), &spaceAPI)
	if err != nil {
		panic(fmt.Sprintf("SPACEAPI is not valid JSON: %v", err))
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	resp, err := http.Get(config.PRESENCEAPI)
	if err == nil {
		defer resp.Body.Close()
		var status Status
		err = json.NewDecoder(resp.Body).Decode(&status)
		if err == nil {
			spaceAPI.State.Open = &status.FuzIsOpen
			if status.FuzIsOpen {
				spaceAPI.State.LastChange = status.LastOpened.Unix()
			} else {
				spaceAPI.State.LastChange = status.LastClosed.Unix()
			}
		}
	} else {
		spaceAPI.State.Open = nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	e := json.NewEncoder(w)
	e.SetIndent("", "    ")
	e.Encode(spaceAPI)
}

func main() {
	http.HandleFunc("/", rootHandler)
	log.Fatal(http.ListenAndServe(":"+config.PORT, nil))
}
