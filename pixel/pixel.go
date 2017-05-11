package pixel

import (
	"image"
	"github.com/yu-ichiko/go-psd/section/header"
)

type Pixel interface {
	image.Image
	SetSource(top, left, bottom, right int, src ...[]byte)
}

func NewPixel(h *header.Header, hasAlpha bool) Pixel {
	switch h.ColorMode {
	case header.ColorModeBitmap, header.ColorModeGrayscale:
	case header.ColorModeRGB:
		return newPixelRGB(h.Depth, hasAlpha)
	case header.ColorModeCMYK:
	}
	return nil
}

func newPixelRGB(depth int, hasAlpha bool) Pixel {
	switch depth {
	case 8:
		if hasAlpha {
			return &PixelNRGBA8{}
		}
		return &PixelNRGB8{}
	case 16:
		if hasAlpha {
			return &PixelNRGBA16{}
		}
		return &PixelNRGB16{}
	//case 32:
	//	if hasAlpha {
	//		return &PixelNRGBA32{}
	//	}
	//	return &PixelNRGB32{}
	}
	return nil
}
