package pixel

import (
	"image"
)

type Pixel interface {
	image.Image
	SetSource(rect image.Rectangle, src ...[]byte)
}

func New(colorMode int, depth int, hasAlpha bool) Pixel {
	switch colorMode {
	case 0, 1:
	case 3:
		return NewPixelRGB(depth, hasAlpha)
	case 4:
	}
	return nil
}

func NewPixelRGB(depth int, hasAlpha bool) Pixel {
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
