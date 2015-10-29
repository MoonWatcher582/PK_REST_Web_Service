package main

import (
	log "github.com/golang/glog"
	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strconv"
)

//	Server comes with a Mongo session
type Server struct {
	session *mgo.Session
}

type Route struct {
	Name        string
	Method      string
	HandlerFunc http.HandlerFunc
}
type route []Route

type Student struct {
	NetID  string `json:"id" bson:"_id"`
	Name   string `json:"name" bson:"name"`
	Major  string `json:"major" bson:"major"`
	Year   int    `json:"year" bson:"year"`
	Grade  int    `json:"grade" bson:"grade"`
	Rating string `json:"rating" bson:"rating"`
}

/*
 * Gets Collection (table).
 *
 *	DB('dbname').C('tablename')
 */
func (s *Server) Collection() *mgo.Collection {
	return s.session.DB("test").C("students")
}

func (s *Server) CreateStudent(w http.ResponseWriter, r *http.Request) {
	// Unmarshal json into Student.
	student := Student{}
	err := json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		// Malformed json.
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)
		return
	}
	log.Infof("%+v", student)

	// Insert into db.
	err = s.Collection().Insert(student)
	if err != nil {
		// Check for duplicates.
		if mgo.IsDup(err) {
			w.WriteHeader(http.StatusConflict)
		} else {
			// Something went wrong.
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(err)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) ReadStudent(w http.ResponseWriter, r *http.Request) {
	// Take URL Query that might have multiple arguments and distill it into one
	// mapping object called query.
	query := make(map[string]interface{})

	// Create the mappping
	for k, v := range r.URL.Query() {
		// If more than one value for a key then fail.
		if len(v) != 1 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		query[k] = v[0]
	}
	log.Infof("%+v", query)

	// Query db.
	result := Student{}
	err := s.Collection().Find(query).One(&result)
	if err != nil {
		// Check if not found.
		if err == mgo.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
		} else {
			// Something went wrong.
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(err)
		}
		return
	}

	// Output result in json.
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Error(err)
	}
}

func (s *Server) UpdateStudents(w http.ResponseWriter, r *http.Request) {
	var students []Student

	// Get all Students. Find(nil) imposes no restriction on search
	err := s.Collection().Find(nil).All(&students)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)
		return
	}

	// Calculate avg of student grades.
	sum := 0
	// First return val is just array index
	for _, student := range students {
		sum += student.Grade
	}
	avg := sum / len(students)

	// Set grade based on avg.
	for _, student := range students {
		switch {
		case student.Grade > (avg + 10):
			student.Rating = "A"
			break
		case student.Grade > (avg - 10):
			student.Rating = "B"
			break
		case student.Grade > (avg - 20):
			student.Rating = "C"
			break
		default:
			continue
		}

		// Update grade.
		err := s.Collection().UpdateId(student.NetID, student)
		if err != nil {
			// Something went wrong.
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(err)
		}
	}
}

func (s *Server) DeleteStudents(w http.ResponseWriter, r *http.Request) {
	year := 0
	var err error

	// Checks if year exists.
	//	Run the assignment, check if ok
	if val, ok := r.URL.Query()["year"]; ok {
		// Checks that only one year is given.
		if len(val) != 1 {
			w.WriteHeader(http.StatusBadRequest)
			log.Error("Did not receive exactly one year.")
			return
		}
		// Checks if the year is an int.
		if year, err = strconv.Atoi(val[0]); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Errorf("Invalid year received [%s]: %v", val[0], err)
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		log.Error("No year was given.")
		return
	}

	// Remove all students with year less than the given year.
	//	Returns metadata about the change, we later print the number of items removed
	changes, err := s.Collection().RemoveAll(bson.M{"year": bson.M{"$lt": year}})
	if err != nil {
		// Something went wrong.
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)
		return
	}

	fmt.Fprintf(w, "Successfully deleted %d students.\n", changes.Removed)
}

func (s *Server) ListStudents(w http.ResponseWriter, r *http.Request) {
	var students []Student

	// Get all Students.
	err := s.Collection().Find(nil).All(&students)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)
	}

	// Output result in json.
	err = json.NewEncoder(w).Encode(students)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)
	}
}

func main() {
	flag.Parse()

	// Connect to mongodb.
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	s := &Server{session}
	log.Info("Opened Mongo Session\n")

	// Set SafeMode so write errors are checked.
	session.SetSafe(&mgo.Safe{})

	// Create routes.
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok\n")) })
	r.HandleFunc("/Student", s.CreateStudent).Methods("POST")
	r.HandleFunc("/Student/getstudent", s.ReadStudent).Methods("GET")
	r.HandleFunc("/Student", s.UpdateStudents).Methods("UPDATE")
	r.HandleFunc("/Student", s.DeleteStudents).Methods("DELETE")
	r.HandleFunc("/Student/listall", s.ListStudents).Methods("GET")
	// Start the server.
	http.ListenAndServe(":1234", r)
}
