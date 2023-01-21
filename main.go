package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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
	Group    string `json:"group"`    // Group Name
	Hostname string `json:"Hostname"` // HOSTNAME
	HostIP   string `json:"HostIP"`   // HOST IP
	IsAlive  bool   `json:"IsAlive"`  // say if there is connection to this host
	PingData string `json:"PingData"` // data from last ping
}

var Hosts []Host // Declare on list of the hosts

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
// ///////// ROUTES FUNCTIONS
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the HomePage!")
	log.Printf("Endpoint Hit: homePage")
}

// ping Functions
func returnAllHosts(w http.ResponseWriter, r *http.Request) {
	log.Printf("EndpointHit: returnAllHosts")
	json.NewEncoder(w).Encode(Hosts)
}
func returnAllHostsWithPing(w http.ResponseWriter, r *http.Request) {
	log.Printf("EndpointHit: returnAllHostsWithPing")
	for i := range Hosts {
		pinger, err := ping.NewPinger(Hosts[i].HostIP)
		pinger.SetPrivileged(true)
		pinger.Timeout = time.Duration(time.Millisecond * 125)
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
		pinger.Timeout = time.Duration(time.Millisecond * 125)
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

// Utils functions
func refresh(w http.ResponseWriter, r *http.Request) {
	Hosts = loadHostsJson("hosts.json")
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

	for i := range groupHosts {
		pinger, err := ping.NewPinger(groupHosts[i].HostIP)
		pinger.SetPrivileged(true)
		pinger.Timeout = time.Duration(time.Millisecond * 125)
		if err != nil {
			panic(err)
		}
		pinger.Count = 3
		err = pinger.Run()

		if err != nil {
			groupHosts[i].PingData = "Error: " + err.Error()
			w.WriteHeader(500)
		} else {

			stats := pinger.Statistics()
			groupHosts[i].PingData = fmt.Sprintf("Packets Sent: %d, Packets Received: %d, Packet Loss: %f%% , RTT Min/Avg/Max: %v/%v/%v",
				stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss, stats.MinRtt, stats.AvgRtt, stats.MaxRtt)
			if stats.PacketsRecv > 0 {
				groupHosts[i].IsAlive = true
			}

		}
	}
	log.Printf("Got Ping data for group")
	json.NewEncoder(w).Encode(groupHosts)
}

// Define routes
var PORT string = "5000"

func handleRequest() {
	// use mux router to handle routes

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", homePage)
	//Util routes
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

	log.Fatal(http.ListenAndServe(":"+PORT, r))
}

// test functions

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

// //////////////////////////////////////////////////////////////
// MAIN FUNCTION
func main() {

	fmt.Println("Rest API v2.0 - Mux Routers")
	fmt.Println("API CREATED BY AVIV MARK RUNNING ON PORT:" + PORT)
	fmt.Println("--------------------------------")
	//Hosts = loadHostsJson("hosts.json")
	Hosts = Get100Servers()
	handleRequest()
}
