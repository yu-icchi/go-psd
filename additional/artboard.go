package additional

import (
	"errors"
	"github.com/yu-ichiko/go-psd/descriptor"
	"github.com/yu-ichiko/go-psd/util"
)

type Artboard struct {
	Name  string         `json:"name"`
	Color *ArtboardColor `json:"color"`
	Type  int            `json:"type"`
	Rect  *ArtboardRect  `json:"rect"`
}

type ArtboardColor struct {
	Red   float64 `json:"red"`
	Green float64 `json:"green"`
	Blue  float64 `json:"blue"`
}

type ArtboardRect struct {
	Top    float64 `json:"top"`
	Left   float64 `json:"left"`
	Bottom float64 `json:"bottom"`
	Right  float64 `json:"right"`
}

// todo...guideIndeces?

func NewArtboard(buf []byte) (*Artboard, error) {
	reader := util.NewReader(buf)
	version, err := reader.ReadInt32()
	if err != nil {
		return nil, err
	}
	if version != 16 {
		return nil, errors.New("invalid artboard")
	}
	desc, err := descriptor.Parse(reader)
	if err != nil {
		return nil, err
	}
	if desc.Class != "artboard" {
		return nil, errors.New("invalid artboard data")
	}

	artboard := &Artboard{}

	for _, item := range desc.Items {
		switch item.Key {
		case "artboardPresetName":
			if name, ok := item.Value.(descriptor.Text); ok {
				artboard.Name = name.String()
			}
		case "artboardBackgroundType":
			if long, ok := item.Value.(descriptor.Integer); ok {
				artboard.Type = long.Integer()
			}
		case "artboardRect":
			if objc, ok := item.Value.(*descriptor.Descriptor); ok {
				rect := &ArtboardRect{}
				for _, artbRect := range objc.Items {
					if doub, ok := artbRect.Value.(descriptor.Double); ok {
						switch artbRect.Key {
						case "Top":
							rect.Top = doub.Number()
						case "Left":
							rect.Left = doub.Number()
						case "Btom":
							rect.Bottom = doub.Number()
						case "Rght":
							rect.Right = doub.Number()
						}
					}
				}
				artboard.Rect = rect
			}
		case "Clr ":
			if objc, ok := item.Value.(*descriptor.Descriptor); ok {
				color := &ArtboardColor{}
				for _, artbRGB := range objc.Items {
					if doub, ok := artbRGB.Value.(descriptor.Double); ok {
						switch artbRGB.Key {
						case "Rd  ":
							color.Red = doub.Number()
						case "Grn ":
							color.Green = doub.Number()
						case "Bl  ":
							color.Blue = doub.Number()
						}
					}
				}
				artboard.Color = color
			}
		}
	}
	return artboard, nil
}
