package color

import (
	"image/color"
	"math"
)

func fromFloat(v, gamma float64) uint32 {
	x := math.Pow(v, gamma)
	switch {
	case x >= 1:
		return 0xffff
	case x <= 0:
		return 0
	default:
		return uint32(x * 0xffff)
	}
}

func toFloat(v uint32, gamma float64) float64 {
	return math.Pow(float64(v)/0xffff, gamma)
}

type NRGBA128 struct {
	R, G, B, A float32
}

func (c NRGBA128) RGBA() (uint32, uint32, uint32, uint32) {
	const gamma = 1.0 / 2.2
	r := fromFloat(float64(c.R), gamma)
	g := fromFloat(float64(c.G), gamma)
	b := fromFloat(float64(c.B), gamma)
	switch {
	case c.A >= 1:
		return r, g, b, 0xffff
	case c.A <= 0:
		return 0, 0, 0, 0
	}
	a := uint32(c.A * 0xffff)
	r = r * a / 0xffff
	g = g * a / 0xffff
	b = b * a / 0xffff
	return r, g, b, a
}

func newNRGBA128Model(c color.Color) color.Color {
	if _, ok := c.(NRGBA128); ok {
		return c
	}
	r, g, b, a := c.RGBA()
	const gamma = 2.2
	fr := float32(toFloat(r, gamma))
	fg := float32(toFloat(g, gamma))
	fb := float32(toFloat(b, gamma))
	switch {
	case a >= 0xffff:
		return NRGBA128{R: fr, G: fg, B: fb, A: 1}
	case a == 0:
		return NRGBA128{}
	}
	fa := 0xffff / float32(a)
	fr *= fa
	fg *= fa
	fb *= fa
	return NRGBA128{R: fr, G: fg, B: fb, A: float32(a) / 0xffff}
}

var (
	NRGBA128Model = color.ModelFunc(newNRGBA128Model)
)
