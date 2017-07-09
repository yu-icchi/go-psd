package image

import (
	"math"
)

func readUint32(b []byte, offset int) uint32 {
	return uint32(b[offset])<<24 | uint32(b[offset+1])<<16 | uint32(b[offset+2])<<8 | uint32(b[offset+3])
}

func readFloat32(b []byte, offset int) float32 {
	return math.Float32frombits(readUint32(b, offset))
}
