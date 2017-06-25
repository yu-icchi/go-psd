package psd

import "errors"

var (
	ErrColorModeData = errors.New("psd: invalid color mode data")
)

type ColorModeData struct {
	Data []byte
}
