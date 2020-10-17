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
}

type SpaceAPI struct {
	API                 string   `json:"api"`
	Space               string   `json:"space"`
	Logo                string   `json:"logo"`
	URL                 string   `json:"url"`
	Location            Location `json:"location"`
	Contact             Contact  `json:"contact"`
	IssueReportChannels []string `json:"issue_report_channels"`
	State               State    `json:"state"`
	Projects            []string `json:"projects"`
}
type Location struct {
	Address string  `json:"address"`
	Lon     float64 `json:"lon"`
	Lat     float64 `json:"lat"`
}
type Contact struct {
	Email   string `json:"email"`
	IRC     string `json:"irc"`
	ML      string `json:"ml"`
	Twitter string `json:"twitter"`
	Matrix  string `json:"matrix"`
}
type Icon struct {
	Open   string `json:"open"`
	Closed string `json:"closed"`
}
type State struct {
	Icon       Icon   `json:"icon"`
	Open       bool   `json:"open"`
	Message    string `json:"message"`
	LastChange int64  `json:"lastchange"`
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
		API:   "0.13",
		Space: "FUZ",
		Logo:  "https://fuz.re/WWW.FUZ.RE_fichiers/5c02b2a84373a.png",
		URL:   "https://fuz.re/",
		Location: Location{
			Address: "11-15 rue de la Réunion, Paris 75020, FRANCE",
			Lat:     48.85343,
			Lon:     2.40308,
		},
		Contact: Contact{
			ML:      "fuz@fuz.re",
			Twitter: "@fuz_re",
			Matrix:  "https://matrix.to/#/#fuz_general:matrix.fuz.re",
		},
		IssueReportChannels: []string{
			"ml",
			"twitter",
		},
		State: State{
			Icon: Icon{
				Open:   "https://presence.fuz.re/img",
				Closed: "https://presence.fuz.re/img",
			},
			Message: "open under conditions: https://wiki.fuz.re/doku.php?id=map",
		},
		Projects: []string{
			"https://wiki.fuz.re/doku.php?id=projets:fuz:start",
		},
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
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	resp, err := http.Get(config.PRESENCEAPI)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	var status Status
	err = json.NewDecoder(resp.Body).Decode(&status)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "err: bad json from presence API: %v", err)
		return
	}
	spaceAPI.State.Open = status.FuzIsOpen
	if status.FuzIsOpen {
		spaceAPI.State.LastChange = status.LastOpened.Unix()
	} else {
		spaceAPI.State.LastChange = status.LastClosed.Unix()
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
