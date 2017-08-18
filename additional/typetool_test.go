package additional

import (
	"github.com/k0kubun/pp"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
)

func TestParseTypetool(t *testing.T) {
	data, err := ioutil.ReadFile("./testdata/typetool_1")
	require.NoError(t, err)
	typetool, err := NewTypeToolObjectSetting(data)
	require.NoError(t, err)
	pp.Println(typetool)
}
