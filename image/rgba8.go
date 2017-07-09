package image

import (
	"image"
	"image/color"
)

type NRGBA8 struct {
	Rect        image.Rectangle
	R           []byte
	G           []byte
	B           []byte
	A           []byte
	Compression int
}

func (p *NRGBA8) CompressionType() int {
	return p.Compression
}

func (p *NRGBA8) Source(rect image.Rectangle, src ...[]byte) {
	p.Rect = rect
	p.R = src[0]
	p.G = src[1]
	p.B = src[2]
	p.A = src[3]
}

func (p *NRGBA8) ColorModel() color.Model {
	return color.NRGBAModel
}

func (p *NRGBA8) Bounds() image.Rectangle {
	return p.Rect
}

func (p *NRGBA8) At(x, y int) color.Color {
	pos := (y-p.Rect.Min.Y)*p.Rect.Dx() + x - p.Rect.Min.X
	return color.NRGBA{
		R: p.R[pos],
		G: p.G[pos],
		B: p.B[pos],
		A: p.A[pos],
	}
}
