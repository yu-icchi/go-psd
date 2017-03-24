package psd

import (
	"errors"
	"io"
	"bufio"
	"fmt"
)

const (
	// File header section (byte length)
	signatureLen = 4
	versionLen = 2
	reservedLen = 6
	imageChannelsLen = 2
	heightLen = 4
	widthLen = 4
	depthLen = 2
	colorModeLen = 2
)

var (
	ErrHeaderFormat = errors.New("invalid psd:header format")
	ErrHeaderVersion = errors.New("invalid psd:header version")
	ErrHeaderChannels = errors.New("invalid psd:header channels")
	ErrHeaderHeight = errors.New("invalid psd:header height")
	ErrHeaderWidth = errors.New("invalid psd:header width")
	ErrHeaderDepth = errors.New("invalid psd:header depth")
	ErrHeaderColorMode = errors.New("invalid psd:header colorMode")
)

type Header struct {
	Signature string
	Version uint16
	Channels uint16
	Height uint32
	Width uint32
	Depth uint16
	colorMode uint16
}

func (h *Header) ColorMode() string {
	switch h.colorMode {
	case 0:
		return "Bitmap"
	case 1:
		return "Grayscale"
	case 2:
		return "Indexed"
	case 3:
		return "RGB"
	case 4:
		return "CMYK"
	case 7:
		return "Multichannel"
	case 8:
		return "Duotone"
	case 9:
		return "Lab"
	default:
		return ""
	}
}

func validSignature(signature string) bool {
	return signature == "8BPS"
}

func validVersion(version uint16) bool {
	return version == 1 || version == 2
}

func validChannels(channels uint16) bool {
	return 1 <= channels && channels <= 56
}

func validHeight(height uint32, version uint16) bool {
	switch version {
	case 1:
		return 1 <= height && height <= 30000
	case 2:
		return 1 <= height && height <= 300000
	default:
		return false
	}
}

func validWidth(width uint32, version uint16) bool {
	switch version {
	case 1:
		return 1 <= width && width <= 30000
	case 2:
		return 1 <= width && width <= 300000
	default:
		return false
	}
}

func validDepth(depth uint16) bool {
	return depth == 1 || depth == 8 || depth == 16 || depth == 32
}

func validColorMode(colorMode uint16) bool {
	return colorMode == 0 ||
		colorMode == 1 ||
		colorMode == 2 ||
		colorMode == 3 ||
		colorMode == 4 ||
		colorMode == 7 ||
		colorMode == 8 ||
		colorMode == 9
}

func readHeader(r io.Reader) (*Header, error) {
	l := signatureLen +
		versionLen +
		reservedLen +
		imageChannelsLen +
		heightLen +
		widthLen +
		depthLen +
		colorModeLen

	buf := make([]byte, l)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, err
	}

	fmt.Println(buf)

	signature := readString(buf, 0, signatureLen)
	if !validSignature(signature) {
		return nil, ErrHeaderFormat
	}

	version := readUint16(buf, signatureLen)
	if !validVersion(version) {
		return nil, ErrHeaderVersion
	}

	channels := readUint16(buf, signatureLen + versionLen + reservedLen)
	if !validChannels(channels) {
		return nil, ErrHeaderChannels
	}

	height := readUint32(buf, signatureLen + versionLen + reservedLen + imageChannelsLen)
	if !validHeight(height, version) {
		return nil, ErrHeaderHeight
	}

	width := readUint32(buf, signatureLen + versionLen + reservedLen + imageChannelsLen + heightLen)
	if !validWidth(width, version) {
		return nil, ErrHeaderWidth
	}

	depth := readUint16(buf, signatureLen + versionLen + reservedLen + imageChannelsLen + heightLen + widthLen)
	if !validDepth(depth) {
		return nil, ErrHeaderDepth
	}

	colorMode := readUint16(buf, signatureLen + versionLen + reservedLen + imageChannelsLen + heightLen + widthLen + depthLen)
	if !validColorMode(colorMode) {
		return nil, ErrHeaderColorMode
	}

	header := &Header{
		Signature: signature,
		Version: version,
		Channels: channels,
		Height: height,
		Width: width,
		Depth: depth,
		colorMode: colorMode,
	}
	return header, nil
}

func writeHeader(w *bufio.Writer, header *Header) {
	w.Write(byteString(header.Signature))
	w.Write(byteUint16(header.Version))
	w.Write([]byte{0, 0, 0, 0, 0, 0})
	w.Write(byteUint16(header.Channels))
	w.Write(byteUint32(header.Height))
	w.Write(byteUint32(header.Width))
	w.Write(byteUint16(header.Depth))
	w.Write(byteUint16(header.colorMode))
	w.Flush()
}
