package additional

import (
	"github.com/yu-ichiko/go-psd/util"
)

func NewBlendInteriorElements(buf []byte) (bool, error) {
	reader := util.NewReader(buf)
	enabled, err := reader.ReadBoolean()
	if err != nil {
		return false, err
	}
	// padding
	if err := reader.Skip(3); err != nil {
		return false, err
	}
	return enabled, nil
}
