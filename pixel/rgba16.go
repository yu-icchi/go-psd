package pixel

import (
	"image"
	"image/color"
	"github.com/yu-ichiko/go-psd/util"
)

type PixelNRGBA16 struct {
	Rect image.Rectangle
	R    []byte
	G    []byte
	B    []byte
	A    []byte
}

func (p *PixelNRGBA16) SetSource(top, left, bottom, right int, src...[]byte) {
	p.Rect = image.Rect(left, top, right, bottom)
	p.R = src[0]
	p.G = src[1]
	p.B = src[2]
	p.A = src[3]
}

func (p *PixelNRGBA16) ColorModel() color.Model {
	return color.NRGBA64Model
}

func (p *PixelNRGBA16) Bounds() image.Rectangle {
	return p.Rect
}

func (p *PixelNRGBA16) At(x, y int) color.Color {
	pos := ((y-p.Rect.Min.Y)*p.Rect.Dx() + x - p.Rect.Min.X) << 1
	return color.NRGBA64{
		R: util.ReadUint16(p.R, pos),
		G: util.ReadUint16(p.G, pos),
		B: util.ReadUint16(p.B, pos),
		A: util.ReadUint16(p.A, pos),
	}
}
