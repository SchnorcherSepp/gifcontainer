package gifcont

import (
	"reflect"
	"testing"
)

func Test_prepareData(t *testing.T) {

	// nil and empty
	for _, in := range [][]byte{nil, {}} {
		data, err := prepareData(in)
		if err != nil {
			t.Error(err)
		}
		if len(data) != header {
			t.Error("wrong len", len(data))
		}
		if data[0] != 0x01 || data[1] != 0x00 || data[2] != 0x00 || data[3] != 0x00 || data[4] != 0x00 {
			t.Error("wrong header")
		}
	}

	// to big
	_, err := prepareData(make([]byte, MaxDataLength+1))
	if err == nil {
		t.Error("no error")
	}

	// max size
	_, err = prepareData(make([]byte, MaxDataLength))
	if err != nil {
		t.Error(err)
	}
}

func Test_extractData(t *testing.T) {
	data := []byte("Hello World!")

	b, err := prepareData(data)
	if err != nil {
		t.Fatal(err)
	}
	datHead := make([]byte, dimension*dimension)
	copy(datHead, b)

	ret, err := extractData(datHead)
	if err != nil {
		t.Fatal(err)
	}

	// check
	if len(ret) != len(data) {
		t.Fatalf("wrong len: ret=%d, orig=%d", len(ret), len(data))
	}
	if !reflect.DeepEqual(ret, data) {
		t.Fatalf("not equal")
	}
}

func Test_extractData_fail(t *testing.T) {
	data := []byte("Hello World!")

	b, err := prepareData(data)
	if err != nil {
		t.Error(err)
	}
	datHead := make([]byte, dimension*dimension)
	copy(datHead, b)

	// check VERSION
	tmp := datHead[0]
	datHead[0] = 0xfa
	_, err = extractData(datHead)
	if err == nil || err.Error() != "wrong version" {
		t.Error("no error:", err)
	}
	datHead[0] = tmp

	// check LEN
	tmp = datHead[2]
	datHead[2] = 0xff
	_, err = extractData(datHead)
	if err == nil || err.Error() != "wrong length" {
		t.Error("no error:", err)
	}
	datHead[2] = tmp

	// check checksum
	tmp = datHead[10]
	datHead[10] = 0xff
	_, err = extractData(datHead)
	if err == nil || err.Error() != "wrong checksum" {
		t.Error("no error:", err)
	}
	datHead[10] = tmp
}
