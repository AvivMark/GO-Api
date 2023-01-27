package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

// Groups Functions
func getGroupsList() []string {
	var Groups []string
	for _, host := range Hosts {
		inSlice := contains(Groups, host.Group)
		if inSlice == false {
			Groups = append(Groups, host.Group)
		}
	}
	return Groups
}

func findGroupHosts(groupName string) []Host {
	var GroupHosts []Host = []Host{}

	for _, host := range Hosts {
		if host.Group == groupName {

			GroupHosts = append(GroupHosts, host)
		}
	}
	return GroupHosts
}

func getGroups(w http.ResponseWriter, r *http.Request) {
	GroupsList := getGroupsList()
	log.Printf("EndpointHit: getGroups")
	json.NewEncoder(w).Encode(GroupsList)
}

func getGroupHosts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["GroupName"]

	groupsHosts := findGroupHosts(key)
	log.Printf("EndpointHit: getGroupHosts for group: " + key)
	json.NewEncoder(w).Encode(groupsHosts)
}

func getGroupAvailable(w http.ResponseWriter, r *http.Request) {
	log.Printf("EndpointHit: getGroupAvailable")
	vars := mux.Vars(r)
	key := vars["GroupName"]
	groupHosts := findGroupHosts(key)
	var wg sync.WaitGroup
	wg.Wait()
	for i := range groupHosts {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			pingHost(&groupHosts[i])
		}(i)
	}
	log.Printf("Got Ping data for group")
	json.NewEncoder(w).Encode(groupHosts)
}
