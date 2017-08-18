package additional

import (
	"github.com/yu-ichiko/go-psd/util"
)

func NewLayerID(buf []byte) (int, error) {
	reader := util.NewReader(buf)
	return reader.ReadInt()
}
