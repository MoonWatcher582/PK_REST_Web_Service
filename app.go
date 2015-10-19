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

func (s *Server) CreateStudent(w http.ResponseWriter, r *http.Request) {
	student := Student{}
	err := json.NewDecoder(r.Body).Decode(&student)
	log.Infof("%+v", student)
	if err != nil {
		log.Error(err)
		return
	}
	err = s.session.DB("test").C("students").Insert(student)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func ReadStudent(w http.ResponseWriter, r *http.Request)   {}
func UpdateStudent(w http.ResponseWriter, r *http.Request) {}
func DeleteStudent(w http.ResponseWriter, r *http.Request) {}
func ListStudent(w http.ResponseWriter, r *http.Request)   {}

func main() {
	flag.Parse()
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		log.Fatal(err)
	}
	s := &Server{session}

	log.Info("Opened Mongo Session\n")
	defer session.Close()

	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok\n")) })
	r.HandleFunc("/student", s.CreateStudent).Methods("POST")
	r.HandleFunc("/student", ReadStudent).Methods("GET")
	http.ListenAndServe(":8000", r)
}
