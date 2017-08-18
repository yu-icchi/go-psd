package image

import (
	"github.com/yu-ichiko/go-psd/util"
	"image"
	"image/color"
)

type NRGBA16 struct {
	Rect        image.Rectangle
	R           []byte
	G           []byte
	B           []byte
	A           []byte
	Compression int
}

func (p *NRGBA16) CompressionType() int {
	return p.Compression
}

func (p *NRGBA16) Source(rect image.Rectangle, src ...[]byte) {
	p.Rect = rect
	p.R = src[0]
	p.G = src[1]
	p.B = src[2]
	p.A = src[3]
}

func (p *NRGBA16) ColorModel() color.Model {
	return color.NRGBA64Model
}

func (p *NRGBA16) Bounds() image.Rectangle {
	return p.Rect
}

func (p *NRGBA16) At(x, y int) color.Color {
	pos := ((y-p.Rect.Min.Y)*p.Rect.Dx() + x - p.Rect.Min.X) << 1
	return color.NRGBA64{
		R: util.ReadUint16(p.R, pos),
		G: util.ReadUint16(p.G, pos),
		B: util.ReadUint16(p.B, pos),
		A: util.ReadUint16(p.A, pos),
	}
}
