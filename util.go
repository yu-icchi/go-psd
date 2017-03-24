package psd

import (
	"encoding/binary"
)

func readString(buf []byte, offset int, limit int) string {
	return string(buf[offset:limit])
}

func readUint16(buf []byte, offset int) uint16 {
	return binary.BigEndian.Uint16(buf[offset:offset+2])
}

func readUint32(buf []byte, offset int) uint32 {
	return binary.BigEndian.Uint32(buf[offset:offset+4])
}

func byteString(str string) []byte {
	return []byte(str)
}

func byteUint16(n uint16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, n)
	return b
}

func byteUint32(n uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, n)
	return b
}
