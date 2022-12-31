package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Host Type Definition and functions

type Host struct {
	ID       string `json:"ID"`       // ID
	Hostname string `json:"Hostname"` // HOSTNAME
	HostIP   string `json:"HostIP"`   // HOST IP
	IsAlive  bool   `json:"isAlive"`  // say if there is connection to this host
}

var Hosts []Host // Declare on list of the hosts

func preDefinedHosts() {
	Hosts = []Host{
		{"1", "server 1", "192.168.1.11", false},
		{"2", "server 2", "192.168.1.12", false},
	}
}

func loadHostsJson(p string) (hosts []Host) {
	content, err := ioutil.ReadFile(p)

	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	var tmp *[]Host
	err = json.Unmarshal(content, &tmp)

	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	return *tmp
}

// //////////////////////////////////////////////////////////////
// ROUTES FUNCTIONS
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func returnAllHosts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("EndpointHit: returnAllHosts")
	json.NewEncoder(w).Encode(Hosts)
}

func getHost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["ID"]

	for _, host := range Hosts {
		if host.ID == key {
			json.NewEncoder(w).Encode(host)
		}
	}
}

func createHost(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var host Host
	json.Unmarshal(reqBody, &host)
	Hosts = append(Hosts, host)
	fmt.Println("EndpointHit: added Host!")
	json.NewEncoder(w).Encode(host)
}

func deleteHost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["ID"]

	for i, host := range Hosts {
		if host.ID == key {
			Hosts = append(Hosts[:i], Hosts[i+1:]...)
		}
	}
	fmt.Println("EndpointHit: deleted Host!")
}

func updateHost(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var hostToUpdate Host
	json.Unmarshal(reqBody, &hostToUpdate)

	for i, host := range Hosts {
		if host.HostIP == hostToUpdate.HostIP {
			Hosts[i] = hostToUpdate
		}
		if host.Hostname == hostToUpdate.Hostname {
			Hosts[i] = hostToUpdate
		}
		if host.ID == hostToUpdate.ID {
			Hosts[i] = hostToUpdate
		}
	}

	fmt.Println("EndpointHit: Updated Host!")
	json.NewEncoder(w).Encode(Hosts)

}
func refresh(w http.ResponseWriter, r *http.Request) {
	Hosts = loadHostsJson("hosts.json")
	fmt.Println("EndpointHit: Refreshed app!")
}

// Define routes
func handleRequest() {
	// use mux router to handle routes

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", homePage)
	r.HandleFunc("/host", createHost).Methods("POST")
	r.HandleFunc("/refresh", refresh)
	r.HandleFunc("/hosts", returnAllHosts)
	r.HandleFunc("/hostUpdate", updateHost).Methods("PUT")
	r.HandleFunc("/host/{ID}", deleteHost).Methods("DELETE")
	r.HandleFunc("/host/{ID}", getHost)
	log.Fatal(http.ListenAndServe(":5000", r))
}

// //////////////////////////////////////////////////////////////
// MAIN FUNCTION
func main() {

	fmt.Println("Rest API v2.0 - Mux Routers")
	fmt.Println("API CREATED BY AVIV MARK")
	fmt.Println("--------------------------------")
	Hosts = loadHostsJson("hosts.json")
	handleRequest()
}
