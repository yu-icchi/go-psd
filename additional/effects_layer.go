package additional

import (
	"errors"
	"github.com/yu-ichiko/go-psd/util"
)

type EffectsLayer struct {
	CommonState *EffectsLayerCommonState
	DropShadow  *EffectsLayerDropShadowAndInnerShadow
	InnerShadow *EffectsLayerDropShadowAndInnerShadow
	OuterGlow   *EffectsLayerOuterGlow
	InnerGlow   *EffectsLayerInnerGlow
	Bevel       *EffectsLayerBevel
	SolidFill   *EffectsLayerSolidFill
}

type EffectsLayerCommonState struct {
	Visible bool
}

type EffectsLayerDropShadowAndInnerShadow struct {
	Blur          int
	Intensity     int
	Angle         float32
	Distance      int
	Color         [2]uint32
	BlendMode     string
	EffectEnabled bool
	UseAngle      bool
	Opacity       uint8
	NativeColor   [2]uint32
}

type EffectsLayerOuterGlow struct {
	Blur          int
	Intensity     int
	Color         [2]uint32
	BlendMode     string
	EffectEnabled bool
	Opacity       uint8
	NativeColor   [2]uint32
}

type EffectsLayerInnerGlow struct {
	Blur          int
	Intensity     int
	Color         [2]uint32
	BlendMode     string
	EffectEnabled bool
	Opacity       uint8
	Invert        uint8
	NativeColor   [2]uint32
}

type EffectsLayerBevel struct {
	Angle              float32
	StrengthDepth      int
	Blur               int
	HighlightBlendMode string
	ShadowBlendMode    string
	HighlightColor     [2]uint32
	ShadowColor        [2]uint32
	BevelStyle         uint8
	HighlightOpacity   uint8
	ShadowOpacity      uint8
	EffectEnabled      bool
	UseAngle           bool
	UpDown             bool
	RealHighlightColor [2]uint32
	RealShadowColor    [2]uint32
}

type EffectsLayerSolidFill struct {
	BlendMode   string
	Color       [4]uint16 // ARGB
	Opacity     uint8
	Enabled     bool
	NativeColor [4]uint16 // ARGB
}

func NewEffectsLayer(buf []byte) (*EffectsLayer, error) {
	reader := util.NewReader(buf)
	version, err := reader.ReadInt16()
	if err != nil {
		return nil, err
	}
	if version != 0 {
		return nil, errors.New("invalid EffectsLayer version")
	}
	count, err := reader.ReadInt16()
	if err != nil {
		return nil, err
	}
	if count != 6 && count != 7 {
		return nil, errors.New("invalid EffectsLayer count")
	}
	effects := &EffectsLayer{}
	for i := 0; i < int(count); i++ {
		sig, err := reader.ReadString(4)
		if err != nil {
			return nil, err
		}
		if sig != "8BIM" {
			return nil, errors.New("invalid EffectsLayer signature")
		}
		key, err := reader.ReadString(4)
		if err != nil {
			return nil, err
		}
		switch key {
		case "cmnS":
			effects.CommonState, err = parseCommonStateInfo(reader)
		case "dsdw":
			effects.DropShadow, err = parseDropShadowAndInnerShadowInfo(reader)
		case "isdw":
			effects.InnerShadow, err = parseDropShadowAndInnerShadowInfo(reader)
		case "oglw":
			effects.OuterGlow, err = parseOuterGlowInfo(reader)
		case "iglw":
			effects.InnerGlow, err = parseInnerGlowInfo(reader)
		case "bevl":
			effects.Bevel, err = parseBevelInfo(reader)
		case "sofi":
			effects.SolidFill, err = parseSolidFill(reader)
		}
		if err != nil {
			return nil, err
		}
	}
	return effects, nil
}

func parseCommonStateInfo(reader *util.Reader) (*EffectsLayerCommonState, error) {
	commonState := &EffectsLayerCommonState{}
	size, err := reader.ReadInt()
	if err != nil {
		return nil, err
	}
	if size != 7 {
		return nil, errors.New("invalid EffectsLayer common state info size")
	}
	version, err := reader.ReadInt()
	if err != nil {
		return nil, err
	}
	if version != 0 {
		return nil, errors.New("invalid EffectsLayer common state info version")
	}
	visible, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	commonState.Visible = visible > 0
	// Unused: always 0
	if err := reader.Skip(2); err != nil {
		return nil, err
	}
	return commonState, nil
}

func parseDropShadowAndInnerShadowInfo(reader *util.Reader) (*EffectsLayerDropShadowAndInnerShadow, error) {
	size, err := reader.ReadInt()
	if err != nil {
		return nil, err
	}
	if size != 41 && size != 51 {
		return nil, errors.New("invalid the remaining items")
	}
	version, err := reader.ReadInt()
	if err != nil {
		return nil, err
	}
	if version != 0 && version != 2 {
		return nil, errors.New("invalid DropShadowAndInnerShadowInfo version")
	}

	dropShadow := &EffectsLayerDropShadowAndInnerShadow{
		Color:       [2]uint32{},
		NativeColor: [2]uint32{},
	}
	dropShadow.Blur, err = reader.ReadInt()
	if err != nil {
		return nil, err
	}
	dropShadow.Intensity, err = reader.ReadInt()
	if err != nil {
		return nil, err
	}
	dropShadow.Angle, err = reader.ReadFloat32()
	if err != nil {
		return nil, err
	}
	dropShadow.Distance, err = reader.ReadInt()
	if err != nil {
		return nil, err
	}

	if err := reader.Skip(2); err != nil {
		return nil, err
	}
	dropShadow.Color[0], err = reader.ReadUInt32()
	if err != nil {
		return nil, err
	}
	dropShadow.Color[1], err = reader.ReadUInt32()
	if err != nil {
		return nil, err
	}

	sig, err := reader.ReadString(4)
	if err != nil {
		return nil, err
	}
	if sig != "8BIM" {
		return nil, errors.New("invalid DropShadowAndInnerShadowInfo signature")
	}
	dropShadow.BlendMode, err = reader.ReadString(4)
	if err != nil {
		return nil, err
	}
	enabled, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	dropShadow.EffectEnabled = enabled > 0
	use, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	dropShadow.UseAngle = use > 0
	dropShadow.Opacity, err = reader.ReadUInt8()
	if err != nil {
		return nil, err
	}

	if err := reader.Skip(2); err != nil {
		return nil, err
	}
	dropShadow.NativeColor[0], err = reader.ReadUInt32()
	if err != nil {
		return nil, err
	}
	dropShadow.NativeColor[1], err = reader.ReadUInt32()
	if err != nil {
		return nil, err
	}
	return dropShadow, nil
}

func parseOuterGlowInfo(reader *util.Reader) (*EffectsLayerOuterGlow, error) {
	size, err := reader.ReadInt()
	if err != nil {
		return nil, err
	}
	if size != 32 && size != 42 {
		return nil, errors.New("invalid the remaining items")
	}
	version, err := reader.ReadInt()
	if err != nil {
		return nil, err
	}
	if version != 0 && version != 2 {
		return nil, errors.New("invalid OuterGlowInfo version")
	}

	outerGlow := &EffectsLayerOuterGlow{}
	outerGlow.Blur, err = reader.ReadInt()
	if err != nil {
		return nil, err
	}
	outerGlow.Intensity, err = reader.ReadInt()
	if err != nil {
		return nil, err
	}

	if err := reader.Skip(2); err != nil {
		return nil, err
	}
	outerGlow.Color[0], err = reader.ReadUInt32()
	if err != nil {
		return nil, err
	}
	outerGlow.Color[1], err = reader.ReadUInt32()
	if err != nil {
		return nil, err
	}

	sig, err := reader.ReadString(4)
	if err != nil {
		return nil, err
	}
	if sig != "8BIM" {
		return nil, errors.New("invalid OuterGlowInfo signature")
	}
	outerGlow.BlendMode, err = reader.ReadString(4)
	if err != nil {
		return nil, err
	}
	enabled, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	outerGlow.EffectEnabled = enabled > 0
	outerGlow.Opacity, err = reader.ReadUInt8()
	if err != nil {
		return nil, err
	}

	if version == 2 {
		if err := reader.Skip(2); err != nil {
			return nil, err
		}
		outerGlow.NativeColor[0], err = reader.ReadUInt32()
		if err != nil {
			return nil, err
		}
		outerGlow.NativeColor[1], err = reader.ReadUInt32()
		if err != nil {
			return nil, err
		}
	}
	return outerGlow, nil
}

func parseInnerGlowInfo(reader *util.Reader) (*EffectsLayerInnerGlow, error) {
	size, err := reader.ReadInt()
	if err != nil {
		return nil, err
	}
	if size != 33 && size != 43 {
		return nil, errors.New("invalid the remaining items")
	}
	version, err := reader.ReadInt()
	if err != nil {
		return nil, err
	}
	if version != 0 && version != 2 {
		return nil, errors.New("invalid OuterGlowInfo version")
	}

	innerGlow := &EffectsLayerInnerGlow{}
	innerGlow.Blur, err = reader.ReadInt()
	if err != nil {
		return nil, err
	}
	innerGlow.Intensity, err = reader.ReadInt()
	if err != nil {
		return nil, err
	}

	if err := reader.Skip(2); err != nil {
		return nil, err
	}
	innerGlow.Color[0], err = reader.ReadUInt32()
	if err != nil {
		return nil, err
	}
	innerGlow.Color[1], err = reader.ReadUInt32()
	if err != nil {
		return nil, err
	}

	sig, err := reader.ReadString(4)
	if err != nil {
		return nil, err
	}
	if sig != "8BIM" {
		return nil, errors.New("invalid OuterGlowInfo signature")
	}
	innerGlow.BlendMode, err = reader.ReadString(4)
	if err != nil {
		return nil, err
	}
	enabled, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	innerGlow.EffectEnabled = enabled > 0
	innerGlow.Opacity, err = reader.ReadUInt8()
	if err != nil {
		return nil, err
	}
	if version == 2 {
		innerGlow.Invert, err = reader.ReadUInt8()
		if err != nil {
			return nil, err
		}
		if err := reader.Skip(2); err != nil {
			return nil, err
		}
		innerGlow.NativeColor[0], err = reader.ReadUInt32()
		if err != nil {
			return nil, err
		}
		innerGlow.NativeColor[1], err = reader.ReadUInt32()
		if err != nil {
			return nil, err
		}
	}
	return innerGlow, nil
}

func parseBevelInfo(reader *util.Reader) (*EffectsLayerBevel, error) {
	size, err := reader.ReadInt()
	if err != nil {
		return nil, err
	}
	if size != 58 && size != 78 {
		return nil, errors.New("invalid BevelInfo size")
	}
	version, err := reader.ReadInt()
	if err != nil {
		return nil, err
	}
	if version != 0 && version != 2 {
		return nil, errors.New("invalid BevelInfo version")
	}

	bevel := &EffectsLayerBevel{}
	bevel.Angle, err = reader.ReadFloat32()
	if err != nil {
		return nil, err
	}
	bevel.StrengthDepth, err = reader.ReadInt()
	if err != nil {
		return nil, err
	}
	bevel.Blur, err = reader.ReadInt()
	if err != nil {
		return nil, err
	}
	sig, err := reader.ReadString(4)
	if err != nil {
		return nil, err
	}
	if sig != "8BIM" {
		return nil, errors.New("invalid BevelInfo signature")
	}
	bevel.HighlightBlendMode, err = reader.ReadString(4)
	if err != nil {
		return nil, err
	}
	sig, err = reader.ReadString(4)
	if err != nil {
		return nil, err
	}
	if sig != "8BIM" {
		return nil, errors.New("invalid BevelInfo signature")
	}
	bevel.ShadowBlendMode, err = reader.ReadString(4)
	if err != nil {
		return nil, err
	}
	if err := reader.Skip(2); err != nil {
		return nil, err
	}
	bevel.HighlightColor[0], err = reader.ReadUInt32()
	if err != nil {
		return nil, err
	}
	bevel.HighlightColor[1], err = reader.ReadUInt32()
	if err != nil {
		return nil, err
	}
	if err := reader.Skip(2); err != nil {
		return nil, err
	}
	bevel.ShadowColor[0], err = reader.ReadUInt32()
	if err != nil {
		return nil, err
	}
	bevel.ShadowColor[1], err = reader.ReadUInt32()
	if err != nil {
		return nil, err
	}
	bevel.BevelStyle, err = reader.ReadByte()
	if err != nil {
		return nil, err
	}
	bevel.HighlightOpacity, err = reader.ReadByte()
	if err != nil {
		return nil, err
	}
	bevel.ShadowOpacity, err = reader.ReadByte()
	if err != nil {
		return nil, err
	}
	enabled, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	bevel.EffectEnabled = enabled > 0
	use, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	bevel.UseAngle = use > 0
	up, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	bevel.UpDown = up == 0
	if version == 2 {
		if err := reader.Skip(2); err != nil {
			return nil, err
		}
		bevel.RealHighlightColor[0], err = reader.ReadUInt32()
		if err != nil {
			return nil, err
		}
		bevel.RealHighlightColor[1], err = reader.ReadUInt32()
		if err != nil {
			return nil, err
		}
		if err := reader.Skip(2); err != nil {
			return nil, err
		}
		bevel.RealShadowColor[0], err = reader.ReadUInt32()
		if err != nil {
			return nil, err
		}
		bevel.RealShadowColor[1], err = reader.ReadUInt32()
		if err != nil {
			return nil, err
		}
	}
	return bevel, nil
}

func parseSolidFill(reader *util.Reader) (*EffectsLayerSolidFill, error) {
	size, err := reader.ReadInt()
	if err != nil {
		return nil, err
	}
	if size != 34 {
		return nil, errors.New("invalid SolidFill size")
	}
	version, err := reader.ReadInt()
	if err != nil {
		return nil, err
	}
	if version != 2 {
		return nil, errors.New("invalid SolidFill version")
	}
	sig, err := reader.ReadString(4)
	if err != nil {
		return nil, err
	}
	if sig != "8BIM" {
		return nil, errors.New("invalid SolidFill signature")
	}

	solidFill := &EffectsLayerSolidFill{}
	solidFill.BlendMode, err = reader.ReadString(4)
	if err != nil {
		return nil, err
	}
	// ARGB
	if err := reader.Skip(2); err != nil {
		return nil, err
	}
	alpha, err := reader.ReadUInt16()
	if err != nil {
		return nil, err
	}
	solidFill.Color[0] = alpha >> 8
	red, err := reader.ReadUInt16()
	if err != nil {
		return nil, err
	}
	solidFill.Color[1] = red >> 8
	green, err := reader.ReadUInt16()
	if err != nil {
		return nil, err
	}
	solidFill.Color[2] = green >> 8
	blue, err := reader.ReadUInt16()
	if err != nil {
		return nil, err
	}
	solidFill.Color[3] = blue >> 8
	solidFill.Opacity, err = reader.ReadUInt8()
	if err != nil {
		return nil, err
	}
	enabled, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	solidFill.Enabled = enabled > 0
	// ARGB
	reader.Skip(2)
	alpha, err = reader.ReadUInt16()
	if err != nil {
		return nil, err
	}
	solidFill.NativeColor[0] = alpha >> 8
	red, err = reader.ReadUInt16()
	if err != nil {
		return nil, err
	}
	solidFill.NativeColor[1] = red >> 8
	green, err = reader.ReadUInt16()
	if err != nil {
		return nil, err
	}
	solidFill.NativeColor[2] = green >> 8
	blue, err = reader.ReadUInt16()
	if err != nil {
		return nil, err
	}
	solidFill.NativeColor[3] = blue >> 8
	return solidFill, nil
}
