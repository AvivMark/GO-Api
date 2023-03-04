package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// host struct
type Group struct {
	ID          string `json:"ID"`               // ID
	Group       string `json:"groupName"`        // Group Name
	Description string `json:"groupDescription"` // Description
	Hosts       []Host `json:"hosts"`            // Hosts
}

// /////////////////////////////////////////////// Host util functions
func getGroupsFromJson(p string) (groups []Group) {
	content, err := ioutil.ReadFile(p)

	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	var tmp *[]Group
	err = json.Unmarshal(content, &tmp)

	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	return *tmp
}

func pingGroup(Group *Group) {
}

// //////
