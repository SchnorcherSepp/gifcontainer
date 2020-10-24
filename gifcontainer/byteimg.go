package gifcont

import (
	"image"
	"image/color"
)

// dimension defines the length and width of the square image (16M Pixel)
const dimension = 4000 // 4k * 4k = 16M

// compiler check
var _ image.Image = (*_ByteImg)(nil)

// _ByteImg display a byte array as an image (1 byte == 1 pixel)
type _ByteImg struct {
	data []byte
}

// ColorModel returns the Image's color model.
func (img *_ByteImg) ColorModel() color.Model {
	return color.RGBAModel
}

// Bounds returns the domain for which At can return non-zero color.
//    4000 pixel x 4000 pixel = 16M Pixel square
func (img *_ByteImg) Bounds() image.Rectangle {
	return image.Rect(0, 0, dimension, dimension)
}

// At returns the color of the pixel at (x, y).
// Outside of the bounds, At() returns a zero color.
func (img *_ByteImg) At(x, y int) color.Color {
	// calc index
	i := index(x, y)

	// get byte
	var byt byte
	if i < len(img.data) {
		byt = img.data[i] // byte from data
	} else {
		byt = 0x00 // index out of range -> zero color
	}

	// return color
	return byte2color(byt)
}
