package main

import(
	"net/http"
	"log"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func main(){

	r := mux.NewRouter()
	r.HandleFunc("/api/items", GetAllItems).Methods("GET")
	r.HandleFunc("/api/items/{id}", GetItem).Methods("GET")
	r.HandleFunc("/api/items", PostItem).Methods("POST")
	r.HandleFunc("/api/items/{id}", DeleteItem).Methods("DELETE")

	http.ListenAndServe(":8086", r)
}


type Item struct {
	ID    string `json:"id" bson:"_id,omitempty"`
	Firstname string    `json:"firstname"`
	Lastname string    `json:"lastname"`
}

var db *mgo.Database

func init() {
	uri := "mongodb://faizul:faizul0@ds227481.mlab.com:27481/sample-data"

	session, err := mgo.Dial(uri)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	db = session.DB("sample-data")
}

func collection() *mgo.Collection {
	return db.C("test-db")
}

// GetAll returns all items from the database.
func GetAll() ([]Item, error) {
	res := []Item{}

	if err := collection().Find(nil).All(&res); err != nil {
		return nil, err
	}

	return res, nil
}

// GetOne returns a single item from the database.
func GetOne(id string) (*Item, error) {
	res := Item{}

	if err := collection().Find(bson.M{"_id": id}).One(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

// Save inserts an item to the database.
func Save(item Item) error {
	return collection().Insert(item)
}

// Remove deletes an item from the database
func Remove(id string) error {
	return collection().Remove(bson.M{"_id": id})
}


func handleError(err error, message string, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(fmt.Sprintf(message, err)))
}

// GetAllItems returns a list of all database items to the response.
func GetAllItems(w http.ResponseWriter, req *http.Request) {
	rs, err := GetAll()
	if err != nil {
		handleError(err, "Failed to load database items: %v", w)
		return
	}

	bs, err := json.Marshal(rs)
	if err != nil {
		handleError(err, "Failed to load marshal data: %v", w)
		return
	}

	w.Write(bs)
}

// GetItem returns a single database item matching given ID parameter.
func GetItem(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	rs, err := GetOne(id)
	if err != nil {
		handleError(err, "Failed to read database: %v", w)
		return
	}

	bs, err := json.Marshal(rs) 
	if err != nil {
		handleError(err, "Failed to marshal data: %v", w)
		return
	}

	w.Write(bs)
}

// PostItem saves an item (form data) into the database.
func PostItem(w http.ResponseWriter, req *http.Request) {
	ID := req.FormValue("id")
	Firstname := req.FormValue("firstname")
	Lastname := req.FormValue("lastname")

	item := Item{ID: ID, Firstname: Firstname, Lastname: Lastname}

	if err := Save(item); err != nil {
		handleError(err, "Failed to save data: %v", w)
		return
	}

	w.Write([]byte("OK"))
}

// DeleteItem removes a single item (identified by parameter) from the database.
func DeleteItem(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	if err := Remove(id); err != nil {
		handleError(err, "Failed to remove item: %v", w)
		return
	}

	w.Write([]byte("OK"))
}