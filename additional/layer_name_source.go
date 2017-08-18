package additional

import "github.com/yu-ichiko/go-psd/util"

func NewLayerNameSource(buf []byte) (string, error) {
	reader := util.NewReader(buf)
	return reader.ReadString(4)
}
