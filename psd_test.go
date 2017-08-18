package psd

import (
	"github.com/k0kubun/pp"
	"github.com/solovev/gopsd"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_psd(t *testing.T) {
	doc, err := gopsd.ParseFromPath("./testdata/artboard_8.psd")
	require.NoError(t, err)
	for _, layer := range doc.Layers {
		pp.Println(layer)
	}
}
