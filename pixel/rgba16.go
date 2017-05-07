package pixel

import "image"

type PixelNRGBA16 struct {
	Rect image.Rectangle
	R    []byte
	G    []byte
	B    []byte
	A    []byte
}
