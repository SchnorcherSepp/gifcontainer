package gifcont

import (
	"testing"
)

func Test_byte2color_color2byte(t *testing.T) {
	// check all possible bytes
	for b := 0x00; b <= 0xff; b++ {

		// get color
		c := byte2color(uint8(b))

		// get byte
		ret := color2byte(c)

		// check
		if ret != uint8(b) {
			t.Errorf("wrong return: ret=%x, b=%x, c=%v", ret, b, c)
		}
	}
}
