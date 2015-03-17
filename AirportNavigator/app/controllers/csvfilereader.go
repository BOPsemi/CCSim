// File:csvfilereader.go

package controllers

import (
	"encoding/csv"
	"io"
	"os"
)

type CsvFileReader struct {
	IATA   []string
	ICAO   []string
	JpName []string
	EnName []string
}

func (c *CsvFileReader) errorHandler(err error) {
	if err != nil {
		panic(err)
	}
}

func (c *CsvFileReader) Open(filePath string) {

	if filePath != "" {
		// open file
		file, err := os.Open(filePath)
		c.errorHandler(err)
		defer file.Close()

		// setup csv reader
		reader := csv.NewReader(file)

		// make list
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			} else {
				c.errorHandler(err)
			}

			c.IATA = append(c.IATA, record[IATA])
			c.ICAO = append(c.ICAO, record[ICAO])
			c.JpName = append(c.JpName, record[JPname])
			c.EnName = append(c.EnName, record[ENname])
		}
	}
}
