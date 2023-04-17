# volumefin

How to Build:
This has mux package for HTTP request handling and header package for page data parsing.

1. go get github.com/golang/gddo/httputil/header

2. go get github.com/gorilla/mux

After that, 
go build volumefin.go

and run volumefin in one terminal. It will start the go process for handling requests under "/calculate".
After that CURL command line can be used for testing input data.

Goal: To create a simple microservice API that can help us understand and track how a particular person's flight path may be queried. 
The API should accept a request that includes a list of flights, which are defined by a source and destination airport code. 
These flights may not be listed in order and will need to be sorted to find the total flight paths starting and ending airports.


API: 
URL: http://<hostname>:8080/calculate
METHOD: POST
DATA: <>
  
Example Request:
curl -X POST  http://localhost:8080/calculate -H 'Content-Type: application/json' -d '[{"source":"IND","dest":"EWR"},{"source":"SFO","dest":"IND"}]'
  
Response:
{
  "source":"SFO","dest":"EWR"
}


This is the implementation of an API for finding/tracking a person or Object
given the data in multiple pairs of cities.
1. Finds Path
2. Checks for Invalid data.
3. Checks for Null Input data.

Design: Create a data struct "flight" which has [src, dest] array. Input will be of the form of an array of flight struct.
Parse through the data and create a Graph along with an Indegree map which keeps tab of which city is visited. Find the
origin city from the inDegree map and the traverse the graph and create a result path. Start and End city of the path is the
output. There could be errors related to invalid data, loops in the cities, unvisited cities etc. Handle them on the way.

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
