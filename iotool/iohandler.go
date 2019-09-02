package iotool

import (
	"encoding/csv"
	"errors"
	"os"
	"strings"
)

/*
IOHandler :interface of io handler
*/
type IOHandler interface {
	WorkingDir() (path string, folderName string) // get working directy info
	OpenFile(filepath string) (*os.File, error)   // open file
	ReadCSV(file *os.File) ([][]string, error)    // read CSV file
}

type ioHandler struct {
	handler IOHandler
}

/*
NewIOHandler :initializer of io handler
*/
func NewIOHandler() IOHandler {
	return new(ioHandler)
}

/*
Get working directory information, return path and folder name
*/
func (io *ioHandler) WorkingDir() (path string, folderName string) {

	// get working directory info
	workDir, err := os.Getwd()
	if err != nil {
		return "", ""
	}

	// get working folder name
	words := strings.Split(workDir, "/")
	workDirName := words[len(words)-1]

	return workDir, workDirName
}

/*
Open file, and return pointer of file
*/
func (io *ioHandler) OpenFile(filepath string) (*os.File, error) {
	// check input
	if filepath == "" {
		return nil, nil
	}

	// open file
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	// return
	return file, nil
}

/*
Read CSV file from pointed file
*/
func (io *ioHandler) ReadCSV(file *os.File) ([][]string, error) {
	// check input
	if file == nil {
		return nil, errors.New("File pointer is nil")
	}

	// initialize CSV file reader
	csvReader := csv.NewReader(file)

	// read CSV file
	record, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	// close file
	file.Close()

	// return
	return record, err
}
