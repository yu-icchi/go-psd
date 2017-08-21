package additional

import (
	"errors"
	"github.com/k0kubun/pp"
	"github.com/yu-ichiko/go-psd/descriptor"
	"github.com/yu-ichiko/go-psd/util"
)

func NewGradientFill(buf []byte) error {
	reader := util.NewReader(buf)
	version, err := reader.ReadInt32()
	if err != nil {
		return err
	}
	if version != 16 {
		return errors.New("")
	}
	desc, err := descriptor.Parse(reader)
	if err != nil {
		return err
	}
	pp.Println(desc)
	return nil
}
