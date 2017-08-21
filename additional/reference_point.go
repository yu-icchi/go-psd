package additional

import (
	"github.com/yu-ichiko/go-psd/util"
)

type Reference struct {
	Points [2]float64
}

// Key is 'fxrp'
func NewReferencePoint(buf []byte) (*Reference, error) {
	reader := util.NewReader(buf)
	ref := &Reference{}
	p1, err := reader.ReadFloat64()
	if err != nil {
		return nil, err
	}
	ref.Points[0] = p1
	p2, err := reader.ReadFloat64()
	if err != nil {
		return nil, err
	}
	ref.Points[1] = p2
	return ref, nil
}
