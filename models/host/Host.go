package host

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/go-ping/ping"
)

// host struct
type Host struct {
	ID       string `json:"ID"`       // ID
	Group    string `json:"group"`    // Group Name
	Hostname string `json:"Hostname"` // HOSTNAME
	HostIP   string `json:"HostIP"`   // HOST IP
	IsAlive  bool   `json:"IsAlive"`  // say if there is connection to this host
	PingData string `json:"PingData"` // data from last ping
}

// /////////////////////////////////////////////// Host util functions
func getHostsFromJson(p string) (hosts []Host) {
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

func pingHost(host *Host) {
	pinger, err := ping.NewPinger(host.HostIP)
	pinger.SetPrivileged(true)
	pinger.Timeout = time.Duration(time.Millisecond * 125)
	if err != nil {
		host.PingData = "Error: " + err.Error()
		host.IsAlive = false
		return
	}
	pinger.Count = 3
	err = pinger.Run()

	if err != nil {
		host.PingData = "Error: " + err.Error()
		host.IsAlive = false
	} else {
		stats := pinger.Statistics()
		host.PingData = fmt.Sprintf("Packets Sent: %d, Packets Received: %d, Packet Loss: %f%% , RTT Min/Avg/Max: %v/%v/%v",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss, stats.MinRtt, stats.AvgRtt, stats.MaxRtt)
		if stats.PacketsRecv > 0 {
			host.IsAlive = true
		}
	}
}

// //////
