package pixel

import "image"

type PixelNRGBA32 struct {
	Rect image.Rectangle
	R    []byte
	G    []byte
	B    []byte
	A    []byte
}
