package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/arriqaaq/flashdb"
)

type Entry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

var db *flashdb.FlashDB

func main() {
	config := &flashdb.Config{Path: "", EvictionInterval: 10}
	var err error
	db, err = flashdb.New(config)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	http.HandleFunc("/add", AddEntry)
	http.HandleFunc("/define/", GetEntry)
	http.HandleFunc("/remove/", RemoveEntry)

	http.ListenAndServe(":8000", nil)
	fmt.Println("Serveur lancé")
}

func AddEntry(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var entry Entry
	err = json.Unmarshal(body, &entry)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusInternalServerError)
		return
	}

	err = db.Update(func(tx *flashdb.Tx) error {
		return tx.Set(entry.Key, entry.Value)
	})

	if err != nil {
		http.Error(w, "Error adding entry to database", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Entry added successfully")
}

func GetEntry(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[len("/define/"):]

	var value string
	err := db.View(func(tx *flashdb.Tx) error {
		val, err := tx.Get(key)
		if err != nil {
			return err
		}
		value = val
		return nil
	})

	if err != nil {
		http.Error(w, "Error getting entry from database", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, value)
}

func RemoveEntry(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[len("/remove/"):]

	err := db.Update(func(tx *flashdb.Tx) error {
		return tx.Delete(key)
	})

	if err != nil {
		http.Error(w, "Error removing entry from database", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Entry removed successfully")
}
