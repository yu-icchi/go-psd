package resources

import (
	"errors"
	"io"

	"fmt"
	"github.com/yu-ichiko/go-psd/util"
)

const (
	length = 4
)

type ImageResourceBlock struct {
	Signature string
	ID        int
	Name      string
	Data      []byte
}

func validSignature(sig string) bool {
	return sig == "8BIM"
}

func Parse(r io.Reader) ([]ImageResourceBlock, int, error) {
	buf := make([]byte, length)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, 0, err
	}
	size := int(util.ReadUint32(buf, 0))
	if size <= 0 {
		return nil, length, nil
	}

	buf = make([]byte, size)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, length, err
	}

	list := []ImageResourceBlock{}
	var read int
	for read < size {
		block := ImageResourceBlock{}

		block.Signature = util.ReadString(buf, read, read+length)
		if !validSignature(block.Signature) {
			return nil, read + length, errors.New("psd: invalid image resource signature")
		}
		read += length

		block.ID = int(util.ReadUint16(buf, read))
		read += 2
		fmt.Println("==========>", block.ID)

		str, l := util.PascalString(buf, read)
		read += l
		fmt.Println("==========>", str, l)
		block.Name = str

		read += util.AdjustAlign2(l)

		size := int(util.ReadUint32(buf, read))
		fmt.Println("==========>", size)
		read += 4

		block.Data = buf[read : read+size]
		read += size
		read += util.AdjustAlign2(size)

		list = append(list, block)
	}

	return list, read + length, nil
}
