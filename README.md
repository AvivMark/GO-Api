# PINGER-API
Api built in go, use this api to create a self monitoring api 

## Configuration

### API Running Port
PORT=5000

### API hosts file - production
HOSTSFILE=hosts.json <br />
this file exists in the main project folder, the file contains the hosts required to monitor with

## API Routes

* /refresh - loading the hosts again from json file

#### Host Routes 
 - (POST request) /host - route to add new host(within the body of the request 
 - (DELETE request) /host/{ID} - Delete host with the use of ID 
 - (GET request) /host/{ID} - Get host with the use of ID,HostIP or Hostname
 - (GET request) /hostAvailable/{ID} - Get host with the use of ID,HostIP or Hostname with new data about their availability
 - (PUT request) /hostUpdate - update host 
 

#### All hosts Routes 
 - (GET request) /hosts - get the list of all hosts
 - (GET request) /hostsAvailable - get the list of all hosts with new data about their availability

#### Group Routes
 - (GET request) /getGroupHosts/{GroupName} - get all hosts under group specified 
 - (GET request) /getGroupAvailable/{GroupName} - get all hosts after sent ping to hosts under group specified 
 - (GET request) /getGroups - get all the groups names 
