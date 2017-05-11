package layer

import (
	"io"

	"github.com/yu-ichiko/go-psd/section/header"
	"github.com/yu-ichiko/go-psd/util"

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
	Clipping    Clipping
	Flags       byte
	Filter      int

	Image       image.Image
}

func (l Layer) Name() string {
	if l.UnicodeName == "" {
		return l.LegacyName
	}
	return l.UnicodeName
}

func (l Layer) Rect() image.Rectangle {
	return image.Rect(l.Left, l.Top, l.Right, l.Bottom)
}

func (l Layer) Width() int {
	return l.Right - l.Left
}

func (l Layer) Height() int {
	return l.Bottom - l.Top
}

func (l Layer) Visible() bool {
	return l.Flags&2 == 1
}

func (l Layer) Obsolete() bool {
	return l.Flags&4 == 1
}

type Channel struct {
	ID     int
	Length int
}

type Clipping int

func (c Clipping) String() string {
	switch c {
	case 0:
		return "base"
	case 1:
		return "non-base"
	}
	return ""
}

type Mask struct {
	Top                int
	Left               int
	Bottom             int
	Right              int
	DefaultColor       int
	Flags              byte

	MaskParametersFlag int
	UserMaskDensity    int
	UserMaskFeather    float64
	VectorMaskDensity  int
	VectorMaskFeather  float64

	Padding            int

	RealFlags          int
	RealBackground     int
	RealTop            int
	RealLeft           int
	RealBottom         int
	RealRight          int
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
	// fmt.Println("== size:", size)

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
	// fmt.Println("=== grobal layer mask info:", size)

	return layers, 0, nil
}
