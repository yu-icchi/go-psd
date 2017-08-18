package additional

import (
	"errors"
	"github.com/yu-ichiko/go-psd/descriptor"
	"github.com/yu-ichiko/go-psd/util"
)

type Typetool struct {
	Transform *TypetoolTransform
	TextData  *descriptor.Descriptor
	WarpData  *descriptor.Descriptor
	Rect      *TypetoolRect
}

type TypetoolTransform struct {
	XX float64
	XY float64
	YX float64
	YY float64
	TX float64
	TY float64
}

type TypetoolRect struct {
	Left   int
	Top    int
	Right  int
	Bottom int
}

func NewTypeToolObjectSetting(buf []byte) (*Typetool, error) {
	reader := util.NewReader(buf)
	version, err := reader.ReadInt16()
	if err != nil {
		return nil, err
	}
	if version != 1 {
		return nil, errors.New("invalid TypeToolObjectSetting version")
	}

	obj := &Typetool{}
	transform := &TypetoolTransform{}
	if transform.XX, err = reader.ReadFloat64(); err != nil {
		return nil, err
	}
	if transform.XY, err = reader.ReadFloat64(); err != nil {
		return nil, err
	}
	if transform.YX, err = reader.ReadFloat64(); err != nil {
		return nil, err
	}
	if transform.YY, err = reader.ReadFloat64(); err != nil {
		return nil, err
	}
	if transform.TX, err = reader.ReadFloat64(); err != nil {
		return nil, err
	}
	if transform.TY, err = reader.ReadFloat64(); err != nil {
		return nil, err
	}
	obj.Transform = transform

	textVersion, err := reader.ReadInt16()
	if err != nil {
		return nil, err
	}
	if textVersion != 50 {
		return nil, errors.New("invalid TypeToolObjectSetting text version")
	}

	descriptorVersion, err := reader.ReadInt32()
	if err != nil {
		return nil, err
	}
	if descriptorVersion != 16 {
		return nil, errors.New("invalid TypeToolObjectSetting descriptor version")
	}

	if obj.TextData, err = descriptor.Parser(reader); err != nil {
		return nil, err
	}

	warpVersion, err := reader.ReadInt16()
	if err != nil {
		return nil, err
	}
	if warpVersion != 1 {
		return nil, errors.New("invalid TypeToolObjectSetting warp version")
	}

	descriptorVersion, err = reader.ReadInt32()
	if err != nil {
		return nil, err
	}
	if descriptorVersion != 16 {
		return nil, errors.New("invalid TypeToolObjectSetting descriptor version")
	}

	if obj.WarpData, err = descriptor.Parser(reader); err != nil {
		return nil, err
	}

	// FIXME: 4byte * 4
	rect := &TypetoolRect{}
	if rect.Left, err = reader.ReadInt(); err != nil {
		return nil, err
	}
	if rect.Top, err = reader.ReadInt(); err != nil {
		return nil, err
	}
	if rect.Right, err = reader.ReadInt(); err != nil {
		return nil, err
	}
	if rect.Bottom, err = reader.ReadInt(); err != nil {
		return nil, err
	}
	obj.Rect = rect

	return obj, nil
}
