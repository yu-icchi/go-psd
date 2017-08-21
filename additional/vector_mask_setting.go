package additional

import (
	"errors"
	"github.com/yu-ichiko/go-psd/pathresource"
	"github.com/yu-ichiko/go-psd/util"
)

type VectorMask struct {
	Invert  bool
	NotLink bool
	Disable bool
	Paths   []*pathresource.Resource
}

func NewVectorMaskSetting(buf []byte) (*VectorMask, error) {
	size := len(buf)
	reader := util.NewReader(buf)
	version, err := reader.ReadInt()
	if err != nil {
		return nil, err
	}
	if version != 3 {
		return nil, errors.New("invalid VectorMaskSetting version")
	}
	flags, err := reader.ReadBytes(4)
	if err != nil {
		return nil, err
	}
	vector := &VectorMask{}
	vector.Invert = flags[0] > 0
	vector.NotLink = flags[1] > 0
	vector.Disable = flags[2] > 0

	num := (size - 10) / 26
	paths := make([]*pathresource.Resource, 0, num)
	for i := 0; i < num; i++ {
		path, err := pathresource.Parse(reader)
		if err != nil {
			return nil, err
		}
		paths = append(paths, path)
	}
	vector.Paths = paths
	return vector, nil
}
