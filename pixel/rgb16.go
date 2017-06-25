package pixel

import (
	"github.com/yu-ichiko/go-psd/util"
	"image"
	"image/color"
)

type PixelNRGB16 struct {
	Rect image.Rectangle
	R    []byte
	G    []byte
	B    []byte
}

func (p *PixelNRGB16) SetSource(rect image.Rectangle, src ...[]byte) {
	p.Rect = rect
	p.R = src[0]
	p.G = src[1]
	p.B = src[2]
}

func (p *PixelNRGB16) ColorModel() color.Model {
	return color.NRGBA64Model
}

func (p *PixelNRGB16) Bounds() image.Rectangle {
	return p.Rect
}

func (p *PixelNRGB16) At(x, y int) color.Color {
	pos := ((y-p.Rect.Min.Y)*p.Rect.Dx() + x - p.Rect.Min.X) << 1
	return color.NRGBA64{
		R: util.ReadUint16(p.R, pos),
		G: util.ReadUint16(p.G, pos),
		B: util.ReadUint16(p.B, pos),
		A: 0xffff,
	}
}
