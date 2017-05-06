package psd

import (
	"errors"
	"io"
)

const (
	// File header section (byte length)
	signatureLen     = 4
	versionLen       = 2
	reservedLen      = 6
	imageChannelsLen = 2
	heightLen        = 4
	widthLen         = 4
	depthLen         = 2
	colorModeLen     = 2
)

var (
	ErrHeaderFormat    = errors.New("invalid psd:header format")
	ErrHeaderVersion   = errors.New("invalid psd:header version")
	ErrHeaderChannels  = errors.New("invalid psd:header channels")
	ErrHeaderHeight    = errors.New("invalid psd:header height")
	ErrHeaderWidth     = errors.New("invalid psd:header width")
	ErrHeaderDepth     = errors.New("invalid psd:header depth")
	ErrHeaderColorMode = errors.New("invalid psd:header colorMode")
)

func fileHeaderSize() int {
	return signatureLen +
		versionLen +
		reservedLen +
		imageChannelsLen +
		heightLen +
		widthLen +
		depthLen +
		colorModeLen
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

func readHeader(r io.Reader) (*FileHeader, int, error) {
	buf := make([]byte, fileHeaderSize())
	l, err := io.ReadFull(r, buf)
	if err != nil {
		return nil, l, err
	}

	signature := readString(buf, 0, signatureLen)
	if !validSignature(signature) {
		return nil, l, ErrHeaderFormat
	}

	version := readUint16(buf, signatureLen)
	if !validVersion(version) {
		return nil, l, ErrHeaderVersion
	}

	channels := readUint16(buf, signatureLen+versionLen+reservedLen)
	if !validChannels(channels) {
		return nil, l, ErrHeaderChannels
	}

	height := readUint32(buf, signatureLen+versionLen+reservedLen+imageChannelsLen)
	if !validHeight(height, version) {
		return nil, l, ErrHeaderHeight
	}

	width := readUint32(buf, signatureLen+versionLen+reservedLen+imageChannelsLen+heightLen)
	if !validWidth(width, version) {
		return nil, l, ErrHeaderWidth
	}

	depth := readUint16(buf, signatureLen+versionLen+reservedLen+imageChannelsLen+heightLen+widthLen)
	if !validDepth(depth) {
		return nil, l, ErrHeaderDepth
	}

	colorMode := readUint16(buf, signatureLen+versionLen+reservedLen+imageChannelsLen+heightLen+widthLen+depthLen)
	if !validColorMode(colorMode) {
		return nil, l, ErrHeaderColorMode
	}

	header := &FileHeader{
		Signature: signature,
		Version:   int(version),
		Channels:  int(channels),
		Height:    int(height),
		Width:     int(width),
		Depth:     int(depth),
		ColorMode: int(colorMode),
	}
	return header, l, nil
}
