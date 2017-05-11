package colormodedata

import (
	"errors"
	"io"

	"github.com/yu-ichiko/go-psd/section/header"
	"github.com/yu-ichiko/go-psd/util"
)

const (
	length = 4
)

var (
	ErrColorModeData = errors.New("psd: invalid color mode data")
)

// ColorModeData Color Mode Data Section
type ColorModeData []byte

// Parse psd color mode data
func Parse(r io.Reader, h *header.Header) (data ColorModeData, read int, err error) {
	var l int
	buf := make([]byte, length)
	if l, err = io.ReadFull(r, buf); err != nil {
		return
	}

	size := int(util.ReadUint32(buf, read))
	read += l
	if size <= 0 {
		return
	}

	if h.ColorMode == header.ColorModeIndexed && size != 768 {
		err = ErrColorModeData
		return
	}

	data = make([]byte, size)
	if l, err = io.ReadFull(r, buf); err != nil {
		return
	}
	read += l
	return
}
