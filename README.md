# volumefin

How to Build:
This has mux package for HTTP request handling and header package for page data parsing.

1. go get github.com/golang/gddo/httputil/header

2. go get github.com/gorilla/mux

After that, 
go build volumefin.go

and run volumefin in one terminal. It will start the go process for handling requests under "/calculate".
After that CURL can be used for testing input data.

This is the implementation of an API for finding/tracking a person or Object
given the data in multiple pairs of cities.
1. Finds Path
2. Checks for Invalid data.
3. Checks for Null Input data.



Example Usage:

1. Valid Data

pavanmettu:test pmettu$ curl -X POST  http://localhost:8080/calculate -H 'Content-Type: application/json' -d '[{"source":"IND","dest":"EWR"},{"source":"SFO","dest":"ATL"},{"source":"GSO", "dest":"IND"}, {"source":"ATL", "dest":"GSO"}]'
{"source":"SFO","dest":"EWR"}

2. Create a Loop in the Cities table

pavanmettu:test pmettu$ curl -X POST  http://localhost:8080/calculate -H 'Content-Type: application/json' -d '[{"source":"IND","dest":"EWR"},{"source":"SFO","dest":"ATL"},{"source":"GSO", "dest":"IND"}, {"source":"ATL", "dest":"SFO"}]'
ERROR: Unvisited Cities, Input Error


3. Pass NULL Input Data

pavanmettu:test pmettu$ curl -X POST  http://localhost:8080/calculate -H 'Content-Type: application/json' -d '[]'
Error: NULL Data
