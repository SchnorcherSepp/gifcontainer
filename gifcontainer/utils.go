package gifcont

import (
	"encoding/binary"
	"image/color"
	"image/color/palette"
)

// reversePlan9 is the inverse of palette.Plan9.
// This var is initialized by the init() function.
var reversePlan9 map[uint32]byte

// init the reversePlan9 var.
func init() {
	// init reversePlan9
	reversePlan9 = make(map[uint32]byte)
	for i, c := range palette.Plan9 {
		key := key(c)
		reversePlan9[key] = uint8(i)
	}
}

//--------------------------------------------------------------------------------------------------------------------//

// key converts a color in a uint32 key.
// Used to access the reverse palette.Plan9 (@see reversePlan9).
func key(c color.Color) uint32 {
	// check nil
	if c == nil {
		return 0 // NIL -> 0
	}

	// calc key
	r, g, b, a := c.RGBA()
	rgba := []byte{uint8(r), uint8(g), uint8(b), uint8(a)}
	key := binary.BigEndian.Uint32(rgba)

	// return
	return key
}

// index calculate an array index with x and y
func index(x, y int) int {
	return x + y*dimension
}

// color2byte return a byte of a color with reversePlan9.
func color2byte(c color.Color) byte {
	key := key(c)

	b, ok := reversePlan9[key]
	if !ok {
		return 0x00 // error
	}

	return b
}

// byte2color convert a byte to a color with palette.Plan9.
func byte2color(byt byte) color.Color {
	// Plan9 is a 256-color palette
	return palette.Plan9[byt]
}
