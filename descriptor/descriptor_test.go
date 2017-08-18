package descriptor

import (
	"github.com/k0kubun/pp"
	"github.com/stretchr/testify/require"
	"github.com/yu-ichiko/go-psd/util"
	"io/ioutil"
	"testing"
)

func TestParser(t *testing.T) {
	data, err := ioutil.ReadFile("./testdata/descriptor_1")
	require.NoError(t, err)
	desc, err := Parser(util.NewReader(data))
	pp.Println(desc)
}
