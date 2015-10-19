package main

import (
	log "github.com/golang/glog"
	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"

	"encoding/json"
	"flag"
	"net/http"
)

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
 * Gets Collection.
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
	// Take URL Query that might have multiple values and distill it into one
	// key value pair.
	query := make(map[string]interface{})

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

func (s *Server) UpdateStudent(w http.ResponseWriter, r *http.Request) {}
func (s *Server) DeleteStudent(w http.ResponseWriter, r *http.Request) {}

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
	// Set SafeMode so write errors are checked.
	session.SetSafe(&mgo.Safe{})
	s := &Server{session}
	log.Info("Opened Mongo Session\n")
	defer session.Close()

	// Create routes.
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok\n")) })
	r.HandleFunc("/student", s.CreateStudent).Methods("POST")
	r.HandleFunc("/student", s.ReadStudent).Methods("GET")
	r.HandleFunc("/student", s.UpdateStudent).Methods("UPDATE")
	r.HandleFunc("/student", s.DeleteStudent).Methods("DELETE")
	r.HandleFunc("/student/listall", s.ListStudents).Methods("GET")
	// Start the server.
	http.ListenAndServe(":8000", r)
}
