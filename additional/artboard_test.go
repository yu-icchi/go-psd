package additional

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
)

func Test_NewArtboard_1(t *testing.T) {
	data, err := ioutil.ReadFile("./testdata/artb_1")
	require.NoError(t, err)
	artboard, err := NewArtboard(data)
	require.NoError(t, err)
	assert.Equal(t, &Artboard{
		Text: "iPhone 6\x00",
		Type: 1,
		Color: &ArtboardColor{
			Red:   255,
			Green: 255,
			Blue:  255,
		},
		Rect: &ArtboardRect{
			Top:    0,
			Left:   750,
			Bottom: 1334,
			Right:  1500,
		},
	}, artboard)
}

func Test_NewArtboard_2(t *testing.T) {
	data, err := ioutil.ReadFile("./testdata/artb_2")
	require.NoError(t, err)
	artboard, err := NewArtboard(data)
	require.NoError(t, err)
	assert.Equal(t, &Artboard{
		Text: "\x00",
		Type: 1,
		Color: &ArtboardColor{
			Red:   255,
			Green: 255,
			Blue:  255,
		},
		Rect: &ArtboardRect{
			Top:    0,
			Left:   1570,
			Bottom: 1068,
			Right:  2627,
		},
	}, artboard)
}
