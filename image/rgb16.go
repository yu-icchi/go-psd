package image

import (
	"github.com/yu-ichiko/go-psd/util"
	"image"
	"image/color"
)

type NRGB16 struct {
	Rect        image.Rectangle
	R           []byte
	G           []byte
	B           []byte
	Compression int
}

func (p *NRGB16) CompressionType() int {
	return p.Compression
}

func (p *NRGB16) Source(rect image.Rectangle, src ...[]byte) {
	p.Rect = rect
	p.R = src[0]
	p.G = src[1]
	p.B = src[2]
}

func (p *NRGB16) ColorModel() color.Model {
	return color.NRGBA64Model
}

func (p *NRGB16) Bounds() image.Rectangle {
	return p.Rect
}

func (p *NRGB16) At(x, y int) color.Color {
	pos := ((y-p.Rect.Min.Y)*p.Rect.Dx() + x - p.Rect.Min.X) << 1
	return color.NRGBA64{
		R: util.ReadUint16(p.R, pos),
		G: util.ReadUint16(p.G, pos),
		B: util.ReadUint16(p.B, pos),
		A: 0xffff,
	}
}
