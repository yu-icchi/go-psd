package pixel

import "image"

type PixelNRGB32 struct {
	Rect image.Rectangle
	R    []byte
	G    []byte
	B    []byte
}
