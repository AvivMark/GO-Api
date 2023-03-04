package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

var PORT string = "5000"        // APP PORT
var JsonFilePath = "hosts.json" // JSON file path
var DebugMode = false

// ///////// UI ROUTES FUNCTIONS
func homePage(w http.ResponseWriter, r *http.Request) {
}

// ///////// Utils functions

// Function creates 100 demo servers for tests only
func Get100Servers() []Host {
	hostsOverload := []Host{}
	for i := 1; i < 121; i++ {
		numStr := strconv.Itoa(i)
		name := "server-" + numStr
		groupNum := strconv.Itoa(i / 12)
		newHost := Host{
			ID:       numStr,
			Hostname: name,
			HostIP:   "172.17.17." + numStr,
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

	//UI Routes
	// Set up a route for the reverse proxy.
	r.HandleFunc("/", homePage)

	log.Fatal(http.ListenAndServe(":"+PORT, r))
}

// //////////////////////////////////////////////////////////////

// MAIN FUNCTION
func main() {

	// Create Test Mode in run
	if len(os.Args) > 1 {
		debugArg := strings.ToLower(os.Args[1])
		if debugArg == "test" {
			Hosts = Get100Servers()
			DebugMode = true
		}
	} else {
		Hosts = getHostsFromJson(JsonFilePath)
	}

	if DebugMode {
		fmt.Println("--------------Debug MODE----------------------")
	}
	fmt.Println("PingGo - The watchdog for your servers")
	fmt.Println("API CREATED BY AVIV MARK RUNNING ON PORT:" + PORT)
	fmt.Println("----------------------------------------------")

	handleRequest()
}
