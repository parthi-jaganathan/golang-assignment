package passwordHandler

import (
	"errors"
	"log"
	"strconv"
	"sync"
)

// struct that holds the hash value
type hashData struct {
	hash string
}

// sequenceIDSync to be able to handle concurrency using sync.RWMutex
var sequenceIDSync = struct {
	sync.RWMutex
	dataMap map[int]*hashData
}{dataMap: make(map[int]*hashData)}

// passHashSync to be able to handle concurrency using sync.RWMutex
var passHashSync = struct {
	sync.RWMutex
	dataMap map[string]*hashData
}{dataMap: make(map[string]*hashData)}

// GenerateSequenceID return a sequence id for the given request which can be used later on to retrieve the password hash
// This method maintaints 2 maps; one with key as the sequence ID and another with key as the hash password
// Both the map holds the reference to the value avoiding duplicate data to be stored in the password hash map
func GenerateSequenceID(password string) int {

	passHashSync.Lock()
	sequenceIDSync.Lock()

	var newAddr *hashData
	existingAddr, ok := passHashSync.dataMap[password]

	// check if the password hash map contains the password, if no, then create a new reference and call it newAddr
	if !ok {
		newAddr = &hashData{hash: password}
		passHashSync.dataMap[password] = newAddr
	} else {
		newAddr = existingAddr
	}

	// increment the sequence number based on the map length
	sequenceID := len(sequenceIDSync.dataMap) + 1
	sequenceIDSync.dataMap[sequenceID] = newAddr

	passHashSync.Unlock()
	sequenceIDSync.Unlock()

	log.Printf("sequenceId %v; passHash %v", sequenceID, password)
	return sequenceID
}

// GetPasswordHashBySequenceID returns the password hash string using the sequece Id
func GetPasswordHashBySequenceID(sequenceID int) (string, error) {

	if sequenceID == 0 { // sanity check, sequence ID start at 1
		return "", errors.New("Invalid sequenceID " + strconv.Itoa(sequenceID))
	}

	sequenceIDSync.RLock()
	val, ok := sequenceIDSync.dataMap[sequenceID]
	sequenceIDSync.RUnlock()

	if !ok {
		return "", errors.New("Invalid sequenceID " + strconv.Itoa(sequenceID))
	}

	data := &val.hash
	resp := *data
	return resp, nil
}
