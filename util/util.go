package util

import (
	"bytes"
	"encoding/binary"
	"unicode/utf16"
)

func ReadString(buf []byte, offset int, limit int) string {
	return string(buf[offset:limit])
}

func ReadUint8(buf []byte, offset int) uint8 {
	return uint8(buf[offset])
}

func ReadUint16(buf []byte, offset int) uint16 {
	return binary.BigEndian.Uint16(buf[offset : offset+2])
}

func ReadInt16(buf []byte, offset int) int16 {
	return int16(ReadUint16(buf, offset))
}

func ReadUint32(buf []byte, offset int) uint32 {
	return binary.BigEndian.Uint32(buf[offset : offset+4])
}

func ReadInt32(buf []byte, offset int) int32 {
	return int32(ReadUint32(buf, offset))
}

func ReadUint64(buf []byte, offset int) uint64 {
	return binary.BigEndian.Uint64(buf[offset : offset+8])
}

func ReadUint64l(buf []byte, offset int) uint64 {
	return binary.LittleEndian.Uint64(buf[offset:offset+8])
}

func ReadClassID(buf []byte) (string, int) {
	size := int(ReadUint32(buf, 0))
	if size == 0 {
		size += 4
	}
	str := ReadString(buf, 4, size + 4)
	return str, size + 4
}

func ByteString(str string) []byte {
	return []byte(str)
}

func ByteUint16(n int) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(n))
	return b
}

func ByteUint32(n int) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(n))
	return b
}

func BytePascalString(str string) ([]byte, int, error) {
	n := len(str)
	if n == 0 {
		return []byte{0}, 1, nil
	}

	buf := &bytes.Buffer{}
	if _, err := buf.Write([]byte{byte(n)}); err != nil {
		return nil, n, err
	}
	if _, err := buf.WriteString(str); err != nil {
		return nil, n, err
	}
	return buf.Bytes(), n, nil
}

func PascalString(buf []byte, offset int) (string, int) {
	size := int(buf[offset])
	if size == 0 {
		return "", 1
	}
	offset += 1
	return string(buf[offset : offset+size]), size
}

func AdjustAlign2(offset int) int {
	if offset&1 != 0 {
		return 1
	}
	return 0
}

func UnicodeString(buf []byte) (string, int) {
	read := 4
	size := ReadUint32(buf, 0)
	if size == 0 {
		return "", read
	}
	data := make([]uint16, size)
	for i := range data {
		data[i] = ReadUint16(buf, 4+i<<1)
		read += 2
	}
	return string(utf16.Decode(data)), read
}

func GetSize(isPSB bool) int {
	if isPSB {
		return 8
	}
	return 4
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	if x == 0 {
		return 0
	}
	return x
}
