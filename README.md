# GO-PINGER-API

## Api built in go, serving hosts(ip) exist in json file


### API Running Port
PORT=5000

### API hosts file - production
HOSTSFILE=hosts.json

### API Routes

(POST request) /host - route to add new host(within the body of the request 
/refresh - loading the hosts again from json file
(GET request) /hosts - get the list of all hosts
(GET request) /hostsAvailable - get the list of all hosts with new data about their availability
(PUT request) /hostUpdate - update host 
(DELETE request) /host/{ID} - Delete host with the use of ID 
(GET request) /host/{ID} - Get host with the use of ID,HostIP or Hostname
(GET request) /hostAvailable/{ID} - Get host with the use of ID,HostIP or Hostname with new data about their availability