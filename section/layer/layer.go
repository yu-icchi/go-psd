package layer

import (
	"image"
	"io"

	"github.com/yu-ichiko/go-psd/section/header"
	"github.com/yu-ichiko/go-psd/util"
)

type Layer struct {
	Index int

	LegacyName string

	Top    int
	Left   int
	Bottom int
	Right  int

	Channels  []Channel
	BlendMode BlendMode
	Opacity   int
	Clipping  Clipping
	Flags     byte
	Filter    int

	Mask              *Mask
	BlendingRanges    *BlendingRanges
	AdditionalInfoMap map[string]AdditionalInfo

	Image image.Image
}

func (l Layer) Name() string {
	return l.LegacyName
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

type BlendMode string

func (b BlendMode) String() string {
	switch b {
	case "pass":
		return "pass through"
	case "norm":
		return "normal"
	case "diss":
		return "dissolve"
	case "dark":
		return "darken"
	case "mul ":
		return "multiply"
	case "idiv":
		return "color burn"
	case "lbrn":
		return "linear burn"
	case "dkCl":
		return "darker color"
	case "lite":
		return "lighten"
	case "scrn":
		return "screen"
	case "div ":
		return "color dodge"
	case "lddg":
		return "linear dodge"
	case "lgCl":
		return "lighter color"
	case "over":
		return "overlay"
	case "sLit":
		return "soft light"
	case "hLit":
		return "hard light"
	case "vLit":
		return "vivid light"
	case "lLit":
		return "linear light"
	case "pLit":
		return "pin light"
	case "hMix":
		return "hard mix"
	case "diff":
		return "difference"
	case "smud":
		return "exclusion"
	case "fsub":
		return "subtract"
	case "fdiv":
		return "divide"
	case "hue ":
		return "hue"
	case "sat ":
		return "saturation"
	case "colr":
		return "color"
	case "lum ":
		return "luminosity"
	}
	return ""
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
	Top          int
	Left         int
	Bottom       int
	Right        int
	DefaultColor int
	Flags        byte

	MaskParametersFlag int
	UserMaskDensity    int
	UserMaskFeather    float64
	VectorMaskDensity  int
	VectorMaskFeather  float64

	Padding int

	RealFlags      int
	RealBackground int
	RealTop        int
	RealLeft       int
	RealBottom     int
	RealRight      int
}

type BlendingRanges struct {
	Black     int
	White     int
	DestRange int
	Channels  []BlendingRangesChannel
}

type BlendingRangesChannel struct {
	SourceRange      int
	DestinationRange int
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
