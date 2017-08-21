package additional

import (
	"errors"
	"github.com/yu-ichiko/go-psd/util"
)

type SectionDivider struct {
	Type      int
	BlendMode string
	SubType   int
}

func NewSectionDividerSetting(buf []byte) (*SectionDivider, error) {
	var err error
	size := len(buf)
	sd := &SectionDivider{}
	reader := util.NewReader(buf)
	sd.Type, err = reader.ReadInt()
	if err != nil {
		return nil, err
	}
	if size >= 12 {
		sig, err := reader.ReadString(4)
		if err != nil {
			return nil, err
		}
		if sig != "8BIM" {
			return nil, errors.New("invalid SectionDividerSetting signature")
		}
		sd.BlendMode, err = reader.ReadString(4)
		if err != nil {
			return nil, err
		}
	}
	if size >= 16 {
		sd.SubType, err = reader.ReadInt()
		if err != nil {
			return nil, err
		}
	}
	return sd, nil
}
