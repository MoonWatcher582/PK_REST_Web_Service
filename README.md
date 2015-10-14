PK REST Web Service

Aedan Dispenza | Alison Wong | Eric Bronner | Jason Davis 

This web service is a (MongoDB?)-backed system (via the mgo driver (https://labix.org/mgo)?) for retrieving student grades for a distributed systems course. The web service is written in Go, supports the POST, GET, DELETE, UPDATE, PUT, and LIST operations, and responds with (JSON|XML).

@TODO
*	Set up Mongo and add test data
*	Create handlers
*	Retrieve variables from the input URL
*	Create routes
*	Connect the database and add r/w
*	Set up JSON parsing
*	Documentation for compilation, use, and database initialization
*	Documentation for parsing, storing, and searching the data
