package colormodedata

import (
	"io"

	"github.com/yu-ichiko/go-psd/psd/util"
)

const (
	length = 4
)

// ColorModeData Color Mode Data Section
type ColorModeData []byte

// Parse psd color mode data
func Parse(r io.Reader) (data ColorModeData, read int, err error) {
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

	data = make([]byte, size)
	if l, err = io.ReadFull(r, buf); err != nil {
		return
	}
	read += l
	return
}
