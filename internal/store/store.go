package store

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

//DB is the type for the in-memory DB
type DB struct {
	Items map[string]string
	mutx  *sync.RWMutex
}

//Set is a method which adds a key-value pair
//to the in-memory store
func (db *DB) Set(key, value string) {
	db.mutx.Lock()
	defer db.mutx.Unlock()
	db.Items[key] = value
}

//Get is a method which returns a value for a
//corresponding key passed as argument
func (db *DB) Get(key string) (string, bool) {
	db.mutx.RLock()
	defer db.mutx.RUnlock()
	value, ok := db.Items[key]
	return value, ok
}

//Delete is a method which deletes the key-value pair
//based on the key provided in the argument
func (db *DB) Delete(key string) {
	db.mutx.Lock()
	defer db.mutx.Unlock()
	delete(db.Items, key)
}

//NewDB is a function that creates and return a
//new instance of the in-memory DB
func NewDB() *DB {
	db := &DB{
		Items: make(map[string]string),
		mutx:  new(sync.RWMutex),
	}

	//case where db file does not exist
	if _, err := os.Stat("dbdata"); os.IsNotExist(err) {
		return db
	}

	//case where db file already exists
	raw, err := ioutil.ReadFile("dbdata")
	if err != nil {
		log.Fatal("DB File read error! ", err.Error())
	}
	buffer := bytes.NewBuffer(raw)
	decoder := gob.NewDecoder(buffer)
	err = decoder.Decode(db)
	if err != nil {
		log.Fatal("GOB decode error! ", err.Error())
	}
	return db
}
