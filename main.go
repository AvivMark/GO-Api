package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/go-ping/ping"
	"github.com/gorilla/mux"
)

// //////////////////////////////////////////////////////////////
// Ping Function
func getPingData(IP string) (data ping.Statistics, e error) {
	pinger, err := ping.NewPinger(IP)
	pinger.SetPrivileged(true)
	pinger.Timeout = time.Duration(time.Millisecond * 300)
	if err != nil {
		noData := ping.Statistics{}
		log.Panic("Error: " + err.Error())
		return noData, err
	}

	pinger.Count = 3
	err = pinger.Run()
	stats := ping.Pinger{}
	pinger.OnFinish(stats.Statistics())

	log.Print("Error: " + err.Error())
	return *stats.Statistics(), err

}

// Host Type Definition and functions

type Host struct {
	ID       string `json:"ID"`       // ID
	Hostname string `json:"Hostname"` // HOSTNAME
	HostIP   string `json:"HostIP"`   // HOST IP
	IsAlive  bool   `json:"IsAlive"`  // say if there is connection to this host
	PingData string `json:"PingData"` // data from last ping
}

var Hosts []Host // Declare on list of the hosts

func preDefinedHosts() {
	Hosts = []Host{
		{"1", "server 1", "192.168.1.11", false, ""},
		{"2", "server 2", "192.168.1.12", false, ""},
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

func (h Host) getHostPingData() error {

	pinger, err := ping.NewPinger(h.HostIP)
	pinger.SetPrivileged(true)
	pinger.Timeout = time.Duration(time.Millisecond * 300)
	log.Print("the value hostStats =", pinger)
	if err != nil {
		panic(err)
	}
	pinger.Count = 3
	err = pinger.Run()

	if err != nil {
		h.PingData = "Error: " + err.Error()
	} else {

		pinger.OnFinish = func(stats *ping.Statistics) {
			h.PingData = fmt.Sprintf("Packets Sent: %d, Packets Received: %d, Packet Loss: %f%% , RTT Min/Avg/Max: %v/%v/%v",
				stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss, stats.MinRtt, stats.AvgRtt, stats.MaxRtt)
			if stats.PacketsRecv > 0 {
				h.IsAlive = true
			}
		}

	}
	return err
}

// //////////////////////////////////////////////////////////////
// ROUTES FUNCTIONS
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func returnAllHostsWithPing(w http.ResponseWriter, r *http.Request) {
	log.Printf("EndpointHit: returnAllHostsWithPing")
	for i := range Hosts {
		pinger, err := ping.NewPinger(Hosts[i].HostIP)
		pinger.SetPrivileged(true)
		pinger.Timeout = time.Duration(time.Millisecond * 300)
		if err != nil {
			panic(err)
		}
		pinger.Count = 3
		err = pinger.Run()

		if err != nil {
			Hosts[i].PingData = "Error: " + err.Error()
			w.WriteHeader(500)
		} else {

			stats := pinger.Statistics()
			Hosts[i].PingData = fmt.Sprintf("Packets Sent: %d, Packets Received: %d, Packet Loss: %f%% , RTT Min/Avg/Max: %v/%v/%v",
				stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss, stats.MinRtt, stats.AvgRtt, stats.MaxRtt)
			if stats.PacketsRecv > 0 {
				Hosts[i].IsAlive = true
			}

		}
	}
	log.Printf("EndpointHit: returned Data!")
	json.NewEncoder(w).Encode(Hosts)
}

func returnAllHosts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("EndpointHit: returnAllHosts")
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
		pinger, err := ping.NewPinger(hostFound.HostIP)
		pinger.SetPrivileged(true)
		pinger.Timeout = time.Duration(time.Millisecond * 300)
		if err != nil {
			panic(err)
		}
		pinger.Count = 3
		err = pinger.Run()

		if err != nil {
			hostFound.PingData = "Error: " + err.Error()
			w.WriteHeader(500)
		} else {

			stats := pinger.Statistics()
			hostFound.PingData = fmt.Sprintf("Packets Sent: %d, Packets Received: %d, Packet Loss: %f%% , RTT Min/Avg/Max: %v/%v/%v",
				stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss, stats.MinRtt, stats.AvgRtt, stats.MaxRtt)
			if stats.PacketsRecv > 0 {
				hostFound.IsAlive = true
			}

			log.Printf("EndpointHit: getHostWithPing for host %s", hostFound.Hostname)
			json.NewEncoder(w).Encode(hostFound)
		}
	} else {
		w.WriteHeader(404)
	}

}

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
			Hosts[i].Hostname = hostToUpdate.Hostname
			Hosts[i].IsAlive = hostToUpdate.IsAlive
		}
		if host.Hostname == hostToUpdate.Hostname {
			Hosts[i].HostIP = hostToUpdate.HostIP
			Hosts[i].IsAlive = hostToUpdate.IsAlive
		}
		if host.ID == hostToUpdate.ID {
			Hosts[i].HostIP = hostToUpdate.HostIP
			Hosts[i].Hostname = hostToUpdate.Hostname
			Hosts[i].IsAlive = hostToUpdate.IsAlive
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
	r.HandleFunc("/hostsAvailable", returnAllHostsWithPing)
	r.HandleFunc("/hostUpdate", updateHost).Methods("PUT")
	r.HandleFunc("/host/{ID}", deleteHost).Methods("DELETE")
	r.HandleFunc("/host/{ID}", getHost)
	r.HandleFunc("/hostAvailable/{ID}", getHostWithPing)
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
