package psd

import (
	"errors"
	"io"
)

const (
	imageResourceLen               = 4
	imageResourceBlockSignature    = "8BIM"
	imageResourceBlockSignatureLen = 4
)

func readImageResources(r io.Reader) ([]ImageResourceBlock, int, error) {

	buf := make([]byte, imageResourceLen)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, 0, err
	}
	size := int(readUint32(buf, 0))
	if size <= 0 {
		return nil, 0, nil
	}

	buf = make([]byte, size)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, 0, err
	}

	// image resource blocks
	list := []ImageResourceBlock{}
	read := 0
	for read < size {
		imgRes := ImageResourceBlock{}

		imgRes.Signature = readString(buf, read, read+imageResourceBlockSignatureLen)
		read += imageResourceBlockSignatureLen
		if imgRes.Signature != imageResourceBlockSignature {
			return nil, 0, errors.New("invalid psd:ImageResourcesBlock signature")
		}

		imgRes.ID = int(readUint16(buf, read))
		read += 2

		str, l := readPascalString(buf, read)
		read += l
		imgRes.Name = str

		read += adjustAlign(l)
		size := int(readUint32(buf, read))
		read += 4

		imgRes.Data = buf[read : read+size]
		read += size
		read += adjustAlign(size)

		list = append(list, imgRes)
	}

	return list, read, nil
}

func readPascalString(buf []byte, offset int) (string, int) {
	size := int(buf[offset])
	if size == 0 {
		return "", 1
	}
	offset += 1
	return string(buf[offset : offset+size]), size
}

func adjustAlign(offset int) int {
	if offset&1 != 0 {
		return 1
	}
	return 0
}
