package psd

import (
	"io"
)

const (
	// The length of the following color data.
	colorModeDataLen = 4
)

func readColorModeData(r io.Reader) (ColorModeData, int, error) {
	buf := make([]byte, colorModeDataLen)
	var read int
	l, err := io.ReadFull(r, buf)
	if err != nil {
		return nil, 0, err
	}
	read += l

	size := readUint32(buf, 0)
	if size <= 0 {
		return nil, read, nil
	}

	buf = make([]byte, size)
	if l, err = io.ReadFull(r, buf); err != nil {
		return nil, read, err
	}
	return buf, read, nil
}
