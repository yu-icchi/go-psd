package pixel

import (
	"image"
	"image/color"
)

type PixelNRGB8 struct {
	Rect image.Rectangle
	R    []byte
	G    []byte
	B    []byte
}

func (p *PixelNRGB8) SetSource(top, left, bottom, right int, src ...[]byte) {
	p.Rect = image.Rect(left, top, right, bottom)
	p.R = src[0]
	p.G = src[1]
	p.B = src[2]
}

func (p *PixelNRGB8) ColorModel() color.Model {
	return color.NRGBAModel
}

func (p *PixelNRGB8) Bounds() image.Rectangle {
	return p.Rect
}

func (p *PixelNRGB8) At(x, y int) color.Color {
	pos := (y-p.Rect.Min.Y)*p.Rect.Dx() + x - p.Rect.Min.X
	return color.NRGBA{
		R: p.R[pos],
		G: p.G[pos],
		B: p.B[pos],
		A: 0xff,
	}
}
