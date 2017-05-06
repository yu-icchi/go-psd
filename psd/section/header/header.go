package header

import (
	"errors"
	"io"

	"github.com/yu-ichiko/go-psd/psd/util"
)

// The color mode of the file
const (
	Bitmap       = "Bitmap"
	Grayscale    = "Grayscale"
	Indexed      = "Indexed"
	RGB          = "RGB"
	CMYK         = "CMYK"
	Multichannel = "Multichannel"
	Duotone      = "Duotone"
	Lab          = "Lab"
)

var (
	headerLens = []int{4, 2, 6, 2, 4, 4, 2, 2}
)

// header error
var (
	ErrHeaderFormat    = errors.New("invalid psd:header format")
	ErrHeaderVersion   = errors.New("invalid psd:header version")
	ErrHeaderChannels  = errors.New("invalid psd:header channels")
	ErrHeaderHeight    = errors.New("invalid psd:header height")
	ErrHeaderWidth     = errors.New("invalid psd:header width")
	ErrHeaderDepth     = errors.New("invalid psd:header depth")
	ErrHeaderColorMode = errors.New("invalid psd:header colorMode")
)

// Header File Header Section
type Header struct {
	Signature string
	Version   int
	Channels  int
	Height    int
	Width     int
	Depth     int
	ColorMode ColorMode
}

// IsPSB psb file
func (h *Header) IsPSB() bool {
	return h.Version == 2
}

// ColorMode
type ColorMode int

// Name color mode name
func (c ColorMode) Name() string {
	switch c {
	case 0:
		return Bitmap
	case 1:
		return Grayscale
	case 2:
		return Indexed
	case 3:
		return RGB
	case 4:
		return CMYK
	case 7:
		return Multichannel
	case 8:
		return Duotone
	case 9:
		return Lab
	default:
		return ""
	}
}

func (c ColorMode) Channels() int {
	switch c {
	case 0, 1, 2:
		return 1
	case 3:
		return 3
	case 4:
		return 4
	}
	return -1
}

func headerLenSize() int {
	size := 0
	for _, n := range headerLens {
		size += n
	}
	return size
}

func validSignature(sig string) bool {
	return sig == "8BPS"
}

func validVersion(version int) bool {
	return version == 1 || version == 2
}

func validChannels(channels int) bool {
	return 1 <= channels && channels <= 56
}

func validHeight(height int, isPSB bool) bool {
	if isPSB {
		return 1 <= height && height <= 300000
	}
	return 1 <= height && height <= 30000
}

func validWidth(width int, isPSB bool) bool {
	if isPSB {
		return 1 <= width && width <= 300000
	}
	return 1 <= width && width <= 30000
}

func validDepth(depth int) bool {
	return depth == 1 || depth == 8 || depth == 16 || depth == 32
}

func validColorMode(colorMode int) bool {
	return colorMode == 0 ||
		colorMode == 1 ||
		colorMode == 2 ||
		colorMode == 3 ||
		colorMode == 4 ||
		colorMode == 7 ||
		colorMode == 8 ||
		colorMode == 9
}

// Parse psd header
func Parse(r io.Reader) (*Header, int, error) {
	read := 0
	buf := make([]byte, headerLenSize())
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, read, err
	}

	header := &Header{}

	// Signature
	read += headerLens[0]
	header.Signature = util.ReadString(buf, 0, read)
	if !validSignature(header.Signature) {
		return nil, read, ErrHeaderFormat
	}

	// Version
	header.Version = int(util.ReadUint16(buf, read))
	if !validVersion(header.Version) {
		return nil, read, ErrHeaderVersion
	}
	read += headerLens[1]

	// Reserved: must be zero
	read += headerLens[2]

	// The number of channels in the image
	header.Channels = int(util.ReadUint16(buf, read))
	if !validChannels(header.Channels) {
		return nil, read, ErrHeaderChannels
	}
	read += headerLens[3]

	// The height of the image in pixels
	header.Height = int(util.ReadUint32(buf, read))
	if !validHeight(header.Height, header.IsPSB()) {
		return nil, read, ErrHeaderHeight
	}
	read += headerLens[4]

	// The width of the image in pixels
	header.Width = int(util.ReadUint32(buf, read))
	if !validWidth(header.Width, header.IsPSB()) {
		return nil, read, ErrHeaderWidth
	}
	read += headerLens[5]

	// Depth
	header.Depth = int(util.ReadUint16(buf, read))
	if !validDepth(header.Depth) {
		return nil, read, ErrHeaderDepth
	}
	read += headerLens[6]

	// The color mode of the file
	header.ColorMode = ColorMode(util.ReadUint16(buf, read))
	read += headerLens[7]

	return header, read, nil
}
