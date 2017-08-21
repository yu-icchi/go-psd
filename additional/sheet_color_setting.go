package additional

import (
	"github.com/yu-ichiko/go-psd/util"
)

type SheetColor struct {
	C1 int
	C2 int
}

func NewSheetColorSetting(buf []byte) (*SheetColor, error) {
	reader := util.NewReader(buf)
	c1, err := reader.ReadInt()
	if err != nil {
		return nil, err
	}
	c2, err := reader.ReadInt()
	if err != nil {
		return nil, err
	}
	return &SheetColor{C1: c1, C2: c2}, nil
}
