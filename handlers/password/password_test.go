package passwordHandler

import (
	"log"
	"strconv"
	"testing"
)

func TestGenerateSequenceID(t *testing.T) {
	var total = 10
	for i := 0; i < total; i++ {
		id := GenerateSequenceID("test" + strconv.Itoa(i))
		if id == 0 {
			t.Error("ID Should be greater than 0")
		}
		if id != i+1 {
			t.Error("ID should be sequential")
		}
	}
}

func TestGenerateSequenceIDAndGetHash(t *testing.T) {
	var total = 10
	for i := 0; i < total; i++ {
		GenerateSequenceID("test" + strconv.Itoa(i))
	}

	for i := 1; i < total+1; i++ {
		_, err := GetPasswordHashBySequenceID(i)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestGenerateSequenceIDInvalidID(t *testing.T) {
	GenerateSequenceID("test")
	_, err1 := GetPasswordHashBySequenceID(0)
	if err1 == nil {
		t.Error(err1)
	}

	_, err2 := GetPasswordHashBySequenceID(2)
	if err2 == nil {
		t.Error(err2)
	}
}

func TestGenerateSequenceIDConcurrent(t *testing.T) {
	var c = make(chan bool)
	var total = 1000
	for i := 0; i <= total; i++ {
		go func() {
			id := GenerateSequenceID("test" + strconv.Itoa(i))
			if id == 0 {
				t.Error("ID Should be greater than 0")
			}

			if id == total {
				c <- true
				log.Println("All Done")
			}
		}()
	}

	<-c
	for i := 1; i < total; i++ {
		go func() {
			_, err := GetPasswordHashBySequenceID(i)
			if err != nil {
				t.Error(err)
			}
		}()
	}
}
