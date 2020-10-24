package gifcont

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"errors"
	"image/gif"
	"reflect"
)

// header : version(1) + length(4) + checksum(16)
const header = 1 + 4 + 16

// MaxDataLength is the max content length that can store in a single image.
const MaxDataLength = dimension*dimension - header

//--------------------------------------------------------------------------------------------------------------------//

// prepareData prepares the data. A header is inserted and the maximum size is ensured.
//
//    --------- header ---------
//    [  1 byte : VERSION      ]
//    [  4 byte : DATA LENGTH  ]
//    [ 16 byte : MD5 CHECKSUM ]
//    -------- content ---------
//    [  n byte : DATA         ]
//
func prepareData(data []byte) ([]byte, error) {

	// convert nil data to empty data
	if data == nil {
		data = []byte{}
	}

	// check data length
	if len(data) > MaxDataLength {
		return []byte{}, errors.New("too much data for one image")
	}

	//-------------------------------------------

	// extend data slice for header
	ret := make([]byte, header+len(data))
	copy(ret[header:], data)

	// add HEADER: version
	ret[0] = 0x01 // add 1 byte

	// add HEADER: data length
	binary.BigEndian.PutUint32(ret[1:5], uint32(len(data))) // add 4 byte

	// add HEADER: checksum
	sum := md5.Sum(data)
	copy(ret[5:21], sum[:]) // add 16 byte

	//-------------------------------------------

	// final size check
	if len(ret) > dimension*dimension {
		return []byte{}, errors.New("final size check fail")
	}

	// success: return data with header
	return ret, nil
}

// extractData split 'in' to 'header' and 'data' and return data.
func extractData(in []byte) ([]byte, error) {

	// check size
	if len(in) != dimension*dimension {
		return []byte{}, errors.New("size check fail")
	}

	// extract header
	version := in[0]
	length := in[1:5]
	sum := in[5:21]
	data := in[21:]

	// check version
	if version != 0x01 {
		return []byte{}, errors.New("wrong version")
	}

	// check length (trim data)
	l := binary.BigEndian.Uint32(length)
	if l > MaxDataLength {
		return []byte{}, errors.New("wrong length")
	}
	data = data[:l]

	// checksum
	sum2 := md5.Sum(data)
	if !reflect.DeepEqual(sum, sum2[:]) {
		return []byte{}, errors.New("wrong checksum")
	}

	// return
	return data, nil
}

//--------------------------------------------------------------------------------------------------------------------//

// Bytes2Image the function packs data into an image (max MaxDataLength bytes per image).
func Bytes2Image(data []byte) (img []byte, err error) {

	// add header
	data, err = prepareData(data)
	if err != nil {
		return []byte{}, err
	}

	// init buffer (the image has 16M pixel always)
	buf := bytes.NewBuffer(make([]byte, 0, dimension*dimension*1.5))

	// generate gif image
	err = gif.Encode(buf, &_ByteImg{data: data}, nil)
	img = buf.Bytes()
	return
}

// Image2Bytes extracts hidden bytes from an image.
func Image2Bytes(img []byte) (data []byte, err error) {

	// check input
	if img == nil {
		img = []byte{} // nil bytes
	}

	// read gif
	gifImg, err := gif.Decode(bytes.NewReader(img))
	if err != nil {
		return []byte{}, err
	}

	// extract bytes
	data = make([]byte, dimension*dimension)
	for x := 0; x < dimension; x++ {
		for y := 0; y < dimension; y++ {
			i := index(x, y)     // calc index
			c := gifImg.At(x, y) // get pixel color
			byt := color2byte(c) // get byte from color
			data[i] = byt        // write byte to array
		}
	}

	// get data
	data, err = extractData(data)
	if err != nil {
		return []byte{}, err
	}

	// return
	return data, nil
}
