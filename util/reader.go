package util

import (
	"bytes"
	"encoding/binary"
	"unicode/utf16"
)

type Reader struct {
	buf *bytes.Reader
	pos int64
}

func NewReader(b []byte) *Reader {
	return &Reader{buf: bytes.NewReader(b), pos: 0}
}

func (r *Reader) ReadByte() (byte, error) {
	var value byte
	if err := binary.Read(r.buf, binary.BigEndian, &value); err != nil {
		return 0, err
	}
	r.pos++
	return value, nil
}

func (r *Reader) ReadBytes(num interface{}) ([]byte, error) {
	n := integer(num)
	value := make([]byte, n)
	if err := binary.Read(r.buf, binary.BigEndian, value); err != nil {
		return nil, err
	}
	r.pos += int64(n)
	return value, nil
}

func (r *Reader) ReadString(n int) (string, error) {
	value := make([]byte, n)
	if err := binary.Read(r.buf, binary.BigEndian, &value); err != nil {
		return "", err
	}
	r.pos += int64(n)
	return string(value), nil
}

func (r *Reader) ReadInt16() (int16, error) {
	var value int16
	if err := binary.Read(r.buf, binary.BigEndian, &value); err != nil {
		return 0, err
	}
	return value, nil
}

func (r *Reader) ReadUInt16() (uint16, error) {
	var value uint16
	if err := binary.Read(r.buf, binary.BigEndian, &value); err != nil {
		return 0, err
	}
	r.pos += 2
	return value, nil
}

func (r *Reader) ReadInt24() (int, error) {
	buffer, err := r.ReadBytes(3)
	if err != nil {
		return 0, err
	}
	value := int(buffer[0]) << 16
	value |= int(buffer[1]) << 8
	value |= int(buffer[2])
	return value, nil
}

func (r *Reader) ReadInt() (int, error) {
	n, err := r.ReadInt32()
	if err != nil {
		return 0, err
	}
	return int(n), nil
}

func (r *Reader) ReadInt32() (int32, error) {
	var value int32
	if err := binary.Read(r.buf, binary.BigEndian, &value); err != nil {
		return 0, err
	}
	r.pos += 4
	return value, nil
}

func (r *Reader) ReadUInt32() (uint32, error) {
	var value uint32
	if err := binary.Read(r.buf, binary.BigEndian, &value); err != nil {
		return 0, err
	}
	r.pos += 4
	return value, nil
}

func (r *Reader) ReadInt64() (int64, error) {
	var value int64
	if err := binary.Read(r.buf, binary.BigEndian, &value); err != nil {
		return 0, err
	}
	r.pos += 8
	return value, nil
}

func (r *Reader) ReadUInt64() (uint64, error) {
	var value uint64
	if err := binary.Read(r.buf, binary.BigEndian, &value); err != nil {
		return 0, err
	}
	r.pos += 8
	return value, nil
}

func (r *Reader) ReadFloat32() (float32, error) {
	var value float32
	if err := binary.Read(r.buf, binary.BigEndian, &value); err != nil {
		return 0, err
	}
	r.pos += 4
	return value, nil
}

func (r *Reader) ReadFloat64() (float64, error) {
	var value float64
	if err := binary.Read(r.buf, binary.BigEndian, &value); err != nil {
		return 0, err
	}
	r.pos += 8
	return value, nil
}

func (r *Reader) ReadPascalString() (string, error) {
	var length byte
	if err := binary.Read(r.buf, binary.BigEndian, &length); err != nil {
		return "", err
	}
	if length == 0 {
		length = 1
	}
	r.pos += 1
	return r.ReadString(int(length))
}

func (r *Reader) ReadUnicodeString() (string, error) {
	num, err := r.ReadInt32()
	if err != nil {
		return "", err
	}
	return r.ReadUnicodeStringLen(int(num))
}

func (r *Reader) ReadUnicodeStringLen(num int) (string, error) {
	str := make([]uint16, num)
	for i := range str {
		if err := binary.Read(r.buf, binary.BigEndian, &str[i]); err != nil {
			return "", err
		}
		r.pos += 2
	}
	return string(utf16.Decode(str)), nil
}

func (r *Reader) ReadDynamicString() (string, error) {
	n, err := r.ReadInt32()
	if err != nil {
		return "", err
	}
	length := int(n)
	if length == 0 {
		length = 4
	}
	return r.ReadString(length)
}

func (r *Reader) Skip(n interface{}) error {
	num := integer(n)
	r.pos += int64(num)
	if _, err := r.buf.Seek(int64(num), 1); err != nil {
		return err
	}
	return nil
}

func integer(num interface{}) int {
	switch i := num.(type) {
	case int64:
		return int(i)
	case int32:
		return int(i)
	case int16:
		return int(i)
	case byte:
		return int(i)
	case int:
		return i
	default:
		return 0
	}
}
