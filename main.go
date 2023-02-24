package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var PORT string = "5000"        // APP PORT
var JsonFilePath = "hosts.json" // JSON file path

// ///////// Basic ROUTES FUNCTIONS
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the HomePage!")
}

// ///////// Utils functions

// Function creates 100 demo servers for tests only
func Get100Servers() []Host {
	hostsOverload := []Host{}
	for i := 1; i < 101; i++ {
		numStr := strconv.Itoa(i)
		name := "server-" + numStr
		groupNum := strconv.Itoa(i / 10)
		newHost := Host{
			ID:       numStr,
			Hostname: name,
			HostIP:   "10.0.0." + numStr,
			IsAlive:  false,
			Group:    "Group" + groupNum,
			PingData: "",
		}
		hostsOverload = append(hostsOverload, newHost)
	}
	return hostsOverload
}

// ROUTE TO  RELOAD HOSTS FROM JSON FILE
func refresh(w http.ResponseWriter, r *http.Request) {
	Hosts = getHostsFromJson(JsonFilePath)
	log.Printf("EndpointHit: Refreshed app!")
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

// Routes main Declaration function
func handleRequest() {
	r := mux.NewRouter()

	// Util routes
	r.HandleFunc("/refresh", refresh)

	// Host routes
	r.HandleFunc("/host", createHost).Methods("POST")
	r.HandleFunc("/hostUpdate", updateHost).Methods("PUT")
	r.HandleFunc("/host/{ID}", deleteHost).Methods("DELETE")
	r.HandleFunc("/host/{ID}", getHost)
	r.HandleFunc("/hostAvailable/{ID}", getHostWithPing)

	// Hosts routes
	r.HandleFunc("/hosts", returnAllHosts)
	r.HandleFunc("/hostsAvailable", returnAllHostsWithPing)

	// Groups Routes
	r.HandleFunc("/getGroupHosts/{GroupName}", getGroupHosts)
	r.HandleFunc("/getGroupAvailable/{GroupName}", getGroupAvailable)
	r.HandleFunc("/getGroups", getGroups)

	//Other Routes
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	handler := cors.Default().Handler(r)
	log.Fatal(http.ListenAndServe(":"+PORT, handler))
}

// //////////////////////////////////////////////////////////////

// MAIN FUNCTION
func main() {

	fmt.Println("Rest API v2.0 - Mux Routers")
	fmt.Println("API CREATED BY AVIV MARK RUNNING ON PORT:" + PORT)
	fmt.Println("----------------------------------------------")
	if len(os.Args) > 1 {

		if os.Args[1] == "Test" {
			Hosts = Get100Servers()
		}
	} else {
		Hosts = getHostsFromJson(JsonFilePath)
	}
	handleRequest()
}
