package additional

import "github.com/yu-ichiko/go-psd/util"

type Locked struct {
	Transparency bool
	Composite    bool
	Position     bool
	All          bool
}

func ParseLocked(buf []byte) Locked {
	locked := int(util.ReadUint32(buf, 0))
	transparency := (locked&(0x01<<0)) > 0 || locked == -2147483648
	composite := (locked&(0x01<<1)) > 0 || locked == -2147483648
	position := (locked&(0x01<<2)) > 0 || locked == -2147483648
	return Locked{
		Transparency: transparency,
		Composite:    composite,
		Position:     position,
		All:          transparency && composite && position,
	}
}
