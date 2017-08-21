package additional

import (
	"errors"
	"github.com/yu-ichiko/go-psd/descriptor"
	"github.com/yu-ichiko/go-psd/util"
)

type Metadata struct {
	Key   string
	Copy  byte
	Value float64
}

func NewMetadataSetting(buf []byte) ([]Metadata, error) {
	reader := util.NewReader(buf)
	count, err := reader.ReadInt()
	if err != nil {
		return nil, err
	}

	list := make([]Metadata, 0, count)
	for i := 0; i < count; i++ {
		sig, err := reader.ReadString(4)
		if err != nil {
			return nil, err
		}
		if sig != "8BIM" && sig != "8B64" {
			return nil, errors.New("invalid Metadata signature")
		}
		key, err := reader.ReadString(4)
		if err != nil {
			return nil, err
		}
		copyOnSheetDup, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		// padding
		if err := reader.Skip(3); err != nil {
			return nil, err
		}
		// length
		length, err := reader.ReadInt()
		if err != nil {
			return nil, err
		}

		var value float64
		if length > 0 {
			version, err := reader.ReadInt()
			if err != nil {
				return nil, err
			}
			if version != 16 {
				return nil, errors.New("invalid Metadata version")
			}
			desc, err := descriptor.Parse(reader)
			if err != nil {
				return nil, err
			}
			if desc.Class == "metadata" {
				for _, item := range desc.Items {
					if doub, ok := item.Value.(descriptor.Double); ok {
						value = doub.Number()
					}
				}
			}
		}

		list = append(list, Metadata{
			Key:   key,
			Copy:  copyOnSheetDup,
			Value: value,
		})
	}
	return list, nil
}
