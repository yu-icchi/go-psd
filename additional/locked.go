package additional

import "github.com/yu-ichiko/go-psd/util"

type Locked struct {
	Transparency bool
	Composite    bool
	Position     bool
	All          bool
}

func NewLocked(buf []byte) (*Locked, error) {
	reader := util.NewReader(buf)
	locked, err := reader.ReadInt32()
	if err != nil {
		return nil, err
	}
	transparency := (locked&(0x01<<0)) > 0 || locked == -2147483648
	composite := (locked&(0x01<<1)) > 0 || locked == -2147483648
	position := (locked&(0x01<<2)) > 0 || locked == -2147483648
	return &Locked{
		Transparency: transparency,
		Composite:    composite,
		Position:     position,
		All:          transparency && composite && position,
	}, nil
}
