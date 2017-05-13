package psd

import "errors"

// header error
var (
	ErrHeaderFormat    = errors.New("psd: invalid header format")
	ErrHeaderVersion   = errors.New("psd: invalid header version")
	ErrHeaderChannels  = errors.New("psd: invalid header channels")
	ErrHeaderHeight    = errors.New("psd: invalid header height")
	ErrHeaderWidth     = errors.New("psd: invalid header width")
	ErrHeaderDepth     = errors.New("psd: invalid header depth")
	ErrHeaderColorMode = errors.New("psd: invalid header colorMode")
)

type Header struct {
	Version   int
	Channels  int
	Height    int
	Width     int
	Depth     int
	ColorMode ColorMode
}

func (h *Header) IsPSB() bool {
	return h.Version == 2
}

var (
	ColorModeBitmap       = ColorMode(0)
	ColorModeGrayscale    = ColorMode(1)
	ColorModeIndexed      = ColorMode(2)
	ColorModeRGB          = ColorMode(3)
	ColorModeCMYK         = ColorMode(4)
	ColorModeMultichannel = ColorMode(7)
	ColorModeDuotone      = ColorMode(8)
	ColorModeLab          = ColorMode(9)
)

type ColorMode int

func (c ColorMode) String() string {
	switch c {
	case ColorModeBitmap:
		return "Bitmap"
	case ColorModeGrayscale:
		return "Grayscale"
	case ColorModeIndexed:
		return "Indexed"
	case ColorModeRGB:
		return "RGB"
	case ColorModeCMYK:
		return "CMYK"
	case ColorModeMultichannel:
		return "Multichannel"
	case ColorModeDuotone:
		return "Duotone"
	case ColorModeLab:
		return "Lab"
	}
	return ""
}