package additional

import (
	"errors"
	"github.com/yu-ichiko/go-psd/descriptor"
	"github.com/yu-ichiko/go-psd/util"
)

func NewObjectEffectsLayerInfo(buf []byte) error {
	reader := util.NewReader(buf)
	version, err := reader.ReadInt()
	if err != nil {
		return err
	}
	if version != 0 {
		return errors.New("invalid Object effects version")
	}
	version, err = reader.ReadInt()
	if err != nil {
		return err
	}
	if version != 16 {
		return errors.New("invalid Object effects Descriptor version")
	}
	_, err = descriptor.Parse(reader)
	if err != nil {
		return err
	}
	//pp.Println(desc)
	return nil
}
