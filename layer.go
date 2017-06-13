package psd

import (
	"github.com/yu-ichiko/go-psd/util"
	"image"
)

var (
	layerSig      = []byte("8BIM")
	additionalSig = []byte("8B64")
)

func newLayer() *Layer {
	return &Layer{
		AdditionalInfos: []*AdditionalInfo{},
	}
}

type Layer struct {
	Index int
	ID    int

	LegacyName string
	Name       string

	Rect  image.Rectangle
	Image image.Image

	Channels     []*Channel
	BlendModeKey BlendModeKey
	Opacity      int
	Clipping     Clipping
	Flags        byte
	Filler       byte

	TransparencyProtected bool
	Visible               bool
	Obsolete              bool
	IrrelevantPixelData   bool

	Mask            *Mask
	BlendingRanges  *BlendingRanges
	AdditionalInfos []*AdditionalInfo
}

func (l *Layer) setRect(top, left, bottom, right int) {
	l.Rect = image.Rect(left, top, right, bottom)
}

func (l *Layer) setAdditionalInfo(addInfo *AdditionalInfo) {

	// Layer ID
	if addInfo.Key == "lyid" {
		l.ID = int(util.ReadUint32(addInfo.Data, 0))
	}

	// Unicode layer name
	if addInfo.Key == "luni" {
		l.Name = util.UnicodeString(addInfo.Data)
	}

	l.AdditionalInfos = append(l.AdditionalInfos, addInfo)
}

type Channel struct {
	ID     int
	Length int
}

type BlendModeKey string

func (b BlendModeKey) String() string {
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

type Clipping byte

func (c Clipping) String() string {
	switch c {
	case 0x00:
		return "base"
	case 0x01:
		return "non-base"
	}
	return ""
}

type Mask struct {
	Rect image.Rectangle

	DefaultColor byte
	Flags        byte
	Padding      int

	RealFlags      byte
	RealBackground byte

	RectEnclosingMask image.Rectangle
}

func (m *Mask) setRect(top, left, bottom, right int) {
	m.Rect = image.Rect(left, top, right, bottom)
}

func (m *Mask) setRectEnclosingMask(top, left, bottom, right int) {
	m.RectEnclosingMask = image.Rect(left, top, right, bottom)
}

func newMask() *Mask {
	return &Mask{}
}

type BlendingRanges struct {
	CompositeGrayBlend *BlendingRangesData
	Channels           []*BlendingRangesData
}

type BlendingRangesData struct {
	Source      int
	Destination int
}

func (b *BlendingRanges) addBlendingRangesData(data *BlendingRangesData) {
	b.Channels = append(b.Channels, data)
}

func newBlendingRanges() *BlendingRanges {
	return &BlendingRanges{
		Channels: []*BlendingRangesData{},
	}
}

type AdditionalInfo struct {
	Key  string
	Data []byte
}
