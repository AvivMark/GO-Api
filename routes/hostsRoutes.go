package routes

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

var Hosts []Host //list of the hosts

func returnAllHostsWithPing(w http.ResponseWriter, r *http.Request) {
	log.Printf("EndpointHit: returnAllHostsWithPing")
	var wg sync.WaitGroup
	for i := range Hosts {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			pingHost(&Hosts[i])
		}(i)
	}
	wg.Wait()
	log.Printf("EndpointHit: returned Data!")
	json.NewEncoder(w).Encode(Hosts)
}

func returnAllHosts(w http.ResponseWriter, r *http.Request) {
	log.Printf("EndpointHit: returnAllHosts")
	json.NewEncoder(w).Encode(Hosts)
}

func getHostWithPing(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["ID"]
	hostFound := Host{}
	for _, host := range Hosts {
		if host.ID == key {
			hostFound = host
		}
		if host.HostIP == key {
			hostFound = host
		}
		if host.Hostname == key {
			hostFound = host
		}
	}
	if hostFound.HostIP != "" {
		pingHost(&hostFound)
		log.Printf("EndpointHit: getHostWithPing for host %s", hostFound.Hostname)
		json.NewEncoder(w).Encode(hostFound)
	} else {
		w.WriteHeader(404)
	}

}

// Host List Edit Functions
func getHost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["ID"]
	hostFound := Host{}
	for _, host := range Hosts {
		if host.ID == key {
			hostFound = host
		}
		if host.HostIP == key {
			hostFound = host
		}
		if host.Hostname == key {
			hostFound = host
		}
	}
	if hostFound.HostIP != "" {
		log.Printf("EndpointHit: getHost for host %s", hostFound.Hostname)
		json.NewEncoder(w).Encode(hostFound)
	} else {
		w.WriteHeader(404)
	}

}

func createHost(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var host Host
	json.Unmarshal(reqBody, &host)
	Hosts = append(Hosts, host)
	log.Printf("EndpointHit: added Host!")
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
	log.Printf("EndpointHit: deleted Host!")
}

func updateHost(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var hostToUpdate Host
	json.Unmarshal(reqBody, &hostToUpdate)

	for i, host := range Hosts {
		if host.HostIP == hostToUpdate.HostIP {
			Hosts[i].Hostname = hostToUpdate.Hostname
			Hosts[i].IsAlive = hostToUpdate.IsAlive
			Hosts[i].Group = hostToUpdate.Group
		}
		if host.Hostname == hostToUpdate.Hostname {
			Hosts[i].HostIP = hostToUpdate.HostIP
			Hosts[i].IsAlive = hostToUpdate.IsAlive
			Hosts[i].Group = hostToUpdate.Group
		}
		if host.ID == hostToUpdate.ID {
			Hosts[i].HostIP = hostToUpdate.HostIP
			Hosts[i].Hostname = hostToUpdate.Hostname
			Hosts[i].IsAlive = hostToUpdate.IsAlive
			Hosts[i].Group = hostToUpdate.Group

		}
	}

	log.Printf("EndpointHit: Updated Host!")
	json.NewEncoder(w).Encode(Hosts)

}
