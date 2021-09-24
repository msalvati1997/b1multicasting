# EXAMPLE B-MULTICAST WITH 3 NODES

---
## How to run example
### Open up 3 terminals and run different commands  from all for each of them
### RUN NODE 1 
```bash
go run main.go -port=:8080 -membersPort=:8081,:8082 -multicastId=MyMulticast -delay=1
```
### RUN NODE 2
```bash
go run main.go -port=:8081 -membersPort=:8080,:8082 -multicastId=MyMulticast -delay=1
```
### RUN NODE 3
```bash
go run main.go -port=:8082 -membersPort=:8080,:8081 -multicastId=MyMulticast -delay=4
```

