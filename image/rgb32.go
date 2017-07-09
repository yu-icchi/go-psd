package image

import (
	"image"
	"image/color"

	pixelColor "github.com/yu-ichiko/go-psd/image/color"
)

type NRGB32 struct {
	Rect        image.Rectangle
	R           []byte
	G           []byte
	B           []byte
	Compression int
}

func (p *NRGB32) CompressionType() int {
	return p.Compression
}

func (p *NRGB32) Source(rect image.Rectangle, src ...[]byte) {
	p.Rect = rect
	p.R = src[0]
	p.G = src[1]
	p.B = src[2]
}

func (p *NRGB32) ColorModel() color.Model {
	return pixelColor.NRGBA128Model
}

func (p *NRGB32) Bounds() image.Rectangle {
	return p.Rect
}

func (p *NRGB32) At(x, y int) color.Color {
	pos := ((y-p.Rect.Min.Y)*p.Rect.Dx() + x - p.Rect.Min.X) << 2
	return pixelColor.NRGBA128{
		R: readFloat32(p.R, pos),
		G: readFloat32(p.G, pos),
		B: readFloat32(p.B, pos),
		A: 1.0,
	}
}
