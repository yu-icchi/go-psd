package enginedata

func Decode2() {

}

type decoder2 struct {
	current  interface{}
	keyStack []string
	stack    []interface{}
}
