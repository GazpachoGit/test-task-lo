# Small task REST API service
## Run
- Install go dev env
- Clone the repository
- Navigate to the repository root folder
- Start the app:
```
go run ./cmd/main.go
```
## Usage
There are 3 endpoints. Examples:
```
curl -d '{"name":"123"}' -X POST http://localhost:9000/tasks
curl -d '{"name":"456"}' -X POST http://localhost:9000/tasks
curl -d '{"name":"789"}' -X POST http://localhost:9000/tasks

curl -X GET http://localhost:9000/tasks
curl -X GET http://localhost:9000/tasks?Status=New

curl -X GET http://localhost:9000/tasks/nTsS84
```