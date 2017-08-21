package pathresource

import (
	"github.com/yu-ichiko/go-psd/util"
)

type Resource struct {
	RecordType int16
	Fill       int16
	Bezier     *Bezier
	Preceding  *Preceding
	Anchor     *Anchor
	Leaving    *Leaving
	Clipboard  *Clipboard
}

type Bezier struct {
	Point [3]int16
}

type Preceding struct {
	Vertical   float64
	Horizontal float64
}

type Anchor struct {
	Vertical   float64
	Horizontal float64
}

type Leaving struct {
	Vertical   float64
	Horizontal float64
}

type Clipboard struct {
	Top        float64
	Left       float64
	Bottom     float64
	Right      float64
	Resolution float64
}

func Parse(reader *util.Reader) (*Resource, error) {
	recordType, err := reader.ReadInt16()
	if err != nil {
		return nil, err
	}
	res := &Resource{
		RecordType: recordType,
	}
	switch res.RecordType {
	case 0, 3:
		if res.Bezier, err = readPathRecord(reader); err != nil {
			return nil, err
		}
	case 1, 2, 4, 5:
		if res.Preceding, res.Anchor, res.Leaving, err = readBezierPoint(reader); err != nil {
			return nil, err
		}
	case 7:
		if res.Clipboard, err = readClipboardRecord(reader); err != nil {
			return nil, err
		}
	case 8:
		if res.Fill, err = readInitialFill(reader); err != nil {
			return nil, err
		}
	default:
		reader.Skip(24)
	}
	return res, nil
}

func readPathRecord(reader *util.Reader) (*Bezier, error) {
	var err error
	bezier := &Bezier{
		Point: [3]int16{},
	}
	if bezier.Point[0], err = reader.ReadInt16(); err != nil {
		return nil, err
	}
	if bezier.Point[1], err = reader.ReadInt16(); err != nil {
		return nil, err
	}
	if bezier.Point[2], err = reader.ReadInt16(); err != nil {
		return nil, err
	}
	if err := reader.Skip(18); err != nil {
		return nil, err
	}
	return bezier, nil
}

func readBezierPoint(reader *util.Reader) (*Preceding, *Anchor, *Leaving, error) {
	var err error
	preceding := &Preceding{}
	if preceding.Vertical, err = reader.ReadPathNumber(); err != nil {
		return nil, nil, nil, err
	}
	if preceding.Horizontal, err = reader.ReadPathNumber(); err != nil {
		return nil, nil, nil, err
	}
	anchor := &Anchor{}
	if anchor.Vertical, err = reader.ReadPathNumber(); err != nil {
		return nil, nil, nil, err
	}
	if anchor.Horizontal, err = reader.ReadPathNumber(); err != nil {
		return nil, nil, nil, err
	}
	leaving := &Leaving{}
	if leaving.Vertical, err = reader.ReadPathNumber(); err != nil {
		return nil, nil, nil, err
	}
	if leaving.Horizontal, err = reader.ReadPathNumber(); err != nil {
		return nil, nil, nil, err
	}
	return preceding, anchor, leaving, nil
}

func readClipboardRecord(reader *util.Reader) (*Clipboard, error) {
	var err error
	clip := &Clipboard{}
	if clip.Top, err = reader.ReadPathNumber(); err != nil {
		return nil, err
	}
	if clip.Left, err = reader.ReadPathNumber(); err != nil {
		return nil, err
	}
	if clip.Bottom, err = reader.ReadPathNumber(); err != nil {
		return nil, err
	}
	if clip.Right, err = reader.ReadPathNumber(); err != nil {
		return nil, err
	}
	if clip.Resolution, err = reader.ReadPathNumber(); err != nil {
		return nil, err
	}
	if err := reader.Skip(4); err != nil {
		return nil, err
	}
	return clip, nil
}

func readInitialFill(reader *util.Reader) (int16, error) {
	initialFill, err := reader.ReadInt16()
	if err != nil {
		return 0, err
	}
	if err := reader.Skip(22); err != nil {
		return 0, err
	}
	return initialFill, nil
}
