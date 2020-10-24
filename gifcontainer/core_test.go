package gifcont_test

import (
	"github.com/SchnorcherSepp/gifcontainer/gifcontainer"
	"io/ioutil"
	"math/rand"
	"reflect"
	"testing"
)

var testRandData []byte

func init() {
	testRandData = make([]byte, 2*gifcont.MaxDataLength)
	rnd := rand.New(rand.NewSource(1337))
	rnd.Read(testRandData)
}

//--------------------------------------------------------------------------------------------------------------------//

func Test_Bytes2Image_Image2Bytes(t *testing.T) {

	// tests
	for _, data := range [][]byte{testRandData[:gifcont.MaxDataLength], testRandData[:0], testRandData[:1], testRandData[:1000000]} {

		// to image
		img, err := gifcont.Bytes2Image(data)
		if err != nil {
			t.Fatal(err)
		}

		// back to bytes
		ret, err := gifcont.Image2Bytes(img)
		if err != nil {
			t.Fatal(err)
		}

		// check
		if len(ret) != len(data) {
			t.Errorf("wrong len: ret=%d, orig=%d", len(ret), len(data))
		}
		if !reflect.DeepEqual(ret, data) {
			t.Errorf("not equal")
		}
	}

	// error
	_, err := gifcont.Bytes2Image(testRandData)
	if err == nil {
		t.Fatal("no error")
	}
}

func Test_Image2Bytes_testFile(t *testing.T) {

	// read test file
	img, err := ioutil.ReadFile("../test/test_v1.gif")
	if err != nil {
		t.Fatal(err)
	}

	// extract data
	ret, err := gifcont.Image2Bytes(img)
	if err != nil {
		t.Fatal(err)
	}

	// check
	data := testRandData[:1000000]
	if len(ret) != len(data) {
		t.Errorf("wrong len: ret=%d, orig=%d", len(ret), len(data))
	}
	if !reflect.DeepEqual(ret, data) {
		t.Errorf("not equal")
	}
}
