package additional

import (
	"errors"
	"github.com/yu-ichiko/go-psd/descriptor"
	"github.com/yu-ichiko/go-psd/util"
)

type SolidColor struct {
	Red   float64
	Green float64
	Blue  float64
}

func NewSolidColorSheetSetting(buf []byte) (*SolidColor, error) {
	reader := util.NewReader(buf)
	version, err := reader.ReadInt()
	if err != nil {
		return nil, err
	}
	if version != 16 {
		return nil, errors.New("invalid SolidColorSheetSetting version")
	}
	desc, err := descriptor.Parse(reader)
	if err != nil {
		return nil, err
	}
	color := &SolidColor{}
	for _, item := range desc.Items {
		objc, ok := item.Value.(*descriptor.Descriptor)
		if !ok {
			continue
		}
		for _, item := range objc.Items {
			switch item.Key {
			case "Rd  ":
				doub, ok := item.Value.(descriptor.Double)
				if !ok {
					continue
				}
				color.Red = doub.Number()
			case "Grn ":
				doub, ok := item.Value.(descriptor.Double)
				if !ok {
					continue
				}
				color.Green = doub.Number()
			case "Bl  ":
				doub, ok := item.Value.(descriptor.Double)
				if !ok {
					continue
				}
				color.Blue = doub.Number()
			}
		}
	}
	return color, nil
}
