package image

import (
	"image"
	"image/color"
)

type NRGB8 struct {
	Rect        image.Rectangle
	R           []byte
	G           []byte
	B           []byte
	Compression int
}

func (p *NRGB8) CompressionType() int {
	return p.Compression
}

func (p *NRGB8) Source(rect image.Rectangle, src ...[]byte) {
	p.Rect = rect
	p.R = src[0]
	p.G = src[1]
	p.B = src[2]
}

func (p *NRGB8) ColorModel() color.Model {
	return color.NRGBAModel
}

func (p *NRGB8) Bounds() image.Rectangle {
	return p.Rect
}

func (p *NRGB8) At(x, y int) color.Color {
	pos := (y-p.Rect.Min.Y)*p.Rect.Dx() + x - p.Rect.Min.X
	return color.NRGBA{
		R: p.R[pos],
		G: p.G[pos],
		B: p.B[pos],
		A: 0xff,
	}
}
