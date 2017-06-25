package psd

import "errors"

var (
	imgResSig = []byte("8BIM")

	ErrImageResourceBlock = errors.New("psd: invalid image resource block")
)

const (
	uniqueIdentifierLen = 2
	actualLen           = 4
)

type ImageResourceBlock struct {
	ID   int
	Name string
	Data []byte
}
