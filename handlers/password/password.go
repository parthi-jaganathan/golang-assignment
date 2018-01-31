package passwordHandler

import (
	"errors"
	"log"
	"strconv"
	"sync"
)

type hashData struct {
	hash string
}

// statsSync to be able to handle concurrency
var sequenceIDSync = struct {
	sync.RWMutex
	dataMap map[int]*hashData
}{dataMap: make(map[int]*hashData)}

// statsSync to be able to handle concurrency
var passHashSync = struct {
	sync.RWMutex
	dataMap map[string]*hashData
}{dataMap: make(map[string]*hashData)}

// GenerateSequenceID return a sequence id for the given request which can be used later on to retrieve the password hash
func GenerateSequenceID(password string) int {

	passHashSync.Lock()
	sequenceIDSync.Lock()

	var newAddr *hashData
	existingAddr, ok := passHashSync.dataMap[password]

	if !ok {
		newAddr = &hashData{hash: password}
		passHashSync.dataMap[password] = newAddr
	} else {
		newAddr = existingAddr
	}

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
