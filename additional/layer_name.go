package additional

import (
	"github.com/yu-ichiko/go-psd/util"
)

func NewLayerName(buf []byte) (string, error) {
	reader := util.NewReader(buf)
	return reader.ReadUnicodeString()
}
