package pixel

import (
	"image"
	"image/color"
)

type PixelNRGBA8 struct {
	Rect image.Rectangle
	R    []byte
	G    []byte
	B    []byte
	A    []byte
}

func (p *PixelNRGBA8) SetSource(top, left, bottom, right int, src ...[]byte) {
	p.Rect = image.Rect(left, top, right, bottom)
	p.R = src[0]
	p.G = src[1]
	p.B = src[2]
	p.A = src[3]
}

func (p *PixelNRGBA8) ColorModel() color.Model {
	return color.NRGBAModel
}

func (p *PixelNRGBA8) Bounds() image.Rectangle {
	return p.Rect
}

func (p *PixelNRGBA8) At(x, y int) color.Color {
	pos := (y-p.Rect.Min.Y)*p.Rect.Dx() + x - p.Rect.Min.X
	return color.NRGBA{
		R: p.R[pos],
		G: p.G[pos],
		B: p.B[pos],
		A: p.A[pos],
	}
}
