package util

import (
	"encoding/binary"
	"unicode/utf16"
)

func ReadString(buf []byte, offset int, limit int) string {
	return string(buf[offset:limit])
}

func ReadUint16(buf []byte, offset int) uint16 {
	return binary.BigEndian.Uint16(buf[offset:offset+2])
}

func ReadUint32(buf []byte, offset int) uint32 {
	return binary.BigEndian.Uint32(buf[offset:offset+4])
}

func ReadUint64(buf []byte, offset int) uint64 {
	return binary.BigEndian.Uint64(buf[offset:offset+8])
}

func ReadUint(buf []byte) uint64 {
	switch len(buf) {
	case 2:
		return uint64(ReadUint16(buf, 0))
	case 4:
		return uint64(ReadUint32(buf, 0))
	case 8:
		return ReadUint64(buf, 0)
	default:
		return 0
	}
}

func ByteString(str string) []byte {
	return []byte(str)
}

func ByteUint16(n uint16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, n)
	return b
}

func ByteUint32(n uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, n)
	return b
}

func PascalString(buf []byte, offset int) (string, int) {
	size := int(buf[offset])
	if size == 0 {
		return "", 1
	}
	offset += 1
	return string(buf[offset:offset+size]), size
}

func AdjustAlign2(offset int) int {
	if offset&1 != 0 {
		return 1
	}
	return 0
}

func UnicodeString(buf []byte) string {
	size := ReadUint32(buf, 0)
	if size == 0 {
		return ""
	}
	data := make([]uint16, size)
	for i := range data {
		data[i] = ReadUint16(buf, 4+i<<1)
	}
	return string(utf16.Decode(data))
}

func GetSize(isPSB bool) int {
	if isPSB {
		return 8
	}
	return 4
}
