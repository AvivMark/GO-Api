# GO-PINGER-API

## Api built in go, serving hosts(ip) exist in json file


### API Running Port
PORT=5000

### API hosts file - production
HOSTSFILE=hosts.json

### API Routes

/host - route to add new host(within the body of the request )
/refresh - loading the hosts again from json file
/hosts - get the list of all hosts
/hostUpdate - update host 
(DELETE request) /host/{ID} - Delete host with the use of ID 
(GET request) /host/{ID} - Get host with the use of ID