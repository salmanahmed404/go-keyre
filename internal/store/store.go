package store

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"log"
	"os"
)

//DB is the type for the in-memory DB
type DB struct {
	Record map[string]interface{}
}

//NewDB is a function that creates and return a
//new instance of the in-memory DB
func NewDB() *DB {
	db := &DB{make(map[string]interface{})}
	if _, err := os.Stat("dbdata"); os.IsNotExist(err) {
		return db
	}

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
