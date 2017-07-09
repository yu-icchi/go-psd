package psd

import (
	"image"
	"fmt"

	psdImage "github.com/yu-ichiko/go-psd/image"
)

type Image interface {
	image.Image
	Source(rect image.Rectangle, src ...[]byte)
}

func newImage(colorMode ColorMode, depth, method int, hasAlpha bool) (Image, error) {
	switch colorMode {
	case ColorModeBitmap:
		return newImageRAW()
	case ColorModeGrayScale:
		return newImageGrayScale()
	case ColorModeRGB:
		return newImageRGB(depth, method, hasAlpha)
	case ColorModeCMYK:
		return newImageCMYK()
	}
	return nil, nil
}

func newImageRAW() (Image, error) {
	return nil, nil
}

func newImageGrayScale() (Image, error) {
	return nil, nil
}

func newImageRGB(depth, method int, hasAlpha bool) (Image, error) {
	switch depth {
	case 8:
		if hasAlpha {
			return &psdImage.NRGBA8{Compression: method}, nil
		}
		return &psdImage.NRGB8{Compression: method}, nil
	case 16:
		if hasAlpha {
			return &psdImage.NRGBA16{Compression: method}, nil
		}
		return &psdImage.NRGB16{Compression: method}, nil
	case 32:
		if hasAlpha {
			return &psdImage.NRGBA32{Compression: method}, nil
		}
		return &psdImage.NRGB32{Compression: method}, nil
	}
	return nil, fmt.Errorf("psd-image: invalid RGB depth %d", depth)
}

func newImageCMYK() (Image, error) {
	return nil, nil
}
