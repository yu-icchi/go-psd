package layer

import (
	"io"

	"github.com/yu-ichiko/go-psd/section/header"
	"github.com/yu-ichiko/go-psd/util"

	"fmt"
	"image"
)

type Layer struct {
	Index       int
	LegacyName  string
	UnicodeName string
	Top         int
	Left        int
	Bottom      int
	Right       int
	Channels    []Channel
	BlendMode   string
	Opacity     int
	Clipping    int
	Flags       int
	Filter      int
	Image       image.Image
}

type Channel struct {
	ID     int
	Length int
}

func (l *Layer) Name() string {
	if l.UnicodeName == "" {
		return l.LegacyName
	}
	return l.UnicodeName
}

func (l *Layer) Width() int {
	return l.Right - l.Left
}

func (l *Layer) Height() int {
	return l.Bottom - l.Top
}

func (l *Layer) IsFolderStart() bool {
	return false
}

func (l *Layer) IsFolderEnd() bool {
	return false
}

func Parse(r io.Reader, header *header.Header) ([]Layer, int, error) {
	var l, read int
	var err error

	size := util.GetSize(header.IsPSB())
	buf := make([]byte, size)

	// Length of the layer and mask information section
	if l, err = io.ReadFull(r, buf); err != nil {
		return nil, 0, err
	}
	read += l

	size = int(util.ReadUint(buf))
	if size <= 0 {
		return nil, read, nil
	}
	fmt.Println("== size:", size)

	// Layer info
	layers, l, err := parseInfo(r, header)
	if err != nil {
		return nil, read, err
	}
	read += l

	// Global layer mask info
	buf = make([]byte, 4)
	if l, err = io.ReadFull(r, buf); err != nil {
		return nil, read, err
	}
	read += l
	size = int(util.ReadUint32(buf, 0))
	fmt.Println("=== grobal layer mask info:", size)

	return layers, 0, nil
}
