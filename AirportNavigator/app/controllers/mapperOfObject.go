// File: mapperOfObject.go

package controllers

import (
	"encoding/json"
	"log"
)

type MapperOfObject struct {
}

// --- Error handler ---
func (c MapperOfObject) errorHandler(err error) {
	if err != nil {
		log.Fatal("Error Recieved")
	}
}

// --- Encoder ---
func (c MapperOfObject) encoder(entity interface{}) []byte {
	data, err := json.Marshal(entity)
	c.errorHandler(err)

	return data
}
