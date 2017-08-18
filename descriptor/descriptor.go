package descriptor

import (
	"github.com/yu-ichiko/go-psd/enginedata"
	"github.com/yu-ichiko/go-psd/util"
)

type (
	Descriptor struct {
		Name  string
		Class string
		Items map[string]*Item
	}
	Item struct {
		Key   string
		Type  string
		Value interface{}
	}
	Reference struct{}
	Double    float64
	UnitFloat struct {
		Type  string
		Value float64
	}
	Text       string
	Enumerated struct {
		Type  string
		Value string
	}
	Integer      int
	LargeInteger int64
	Boolean      bool
	Class        struct {
		Name string
		ID   string
	}
	Alias string
)

func (t Text) String() string {
	return string(t)
}

func (d Double) Number() float64 {
	return float64(d)
}

func (i Integer) Integer() int {
	return int(i)
}

func Parser(reader *util.Reader) (*Descriptor, error) {
	var err error

	obj := &Descriptor{}
	obj.Name, err = reader.ReadUnicodeString()
	if err != nil {
		return nil, err
	}
	obj.Class, err = reader.ReadDynamicString()
	if err != nil {
		return nil, err
	}
	num, err := reader.ReadInt()
	if err != nil {
		return nil, err
	}

	items := map[string]*Item{}
	for i := 0; i < num; i++ {
		item, err := parseItem(reader)
		if err != nil {
			return nil, err
		}
		items[item.Key] = item
	}
	obj.Items = items

	return obj, nil
}

func parseItem(reader *util.Reader) (*Item, error) {
	var err error

	item := &Item{}
	item.Key, err = reader.ReadDynamicString()
	if err != nil {
		return nil, err
	}
	item.Type, err = reader.ReadString(4)
	if err != nil {
		return nil, err
	}

	switch item.Type {
	case "obj ":
		size, err := reader.ReadInt()
		if err != nil {
			return nil, err
		}
		for i := 0; i < size; i++ {
			typ, err := reader.ReadString(4)
			if err != nil {
				return nil, err
			}
			switch typ {
			// todo...
			}
		}
	case "Objc", "Glb0":
		item.Value, err = Parser(reader)
		if err != nil {
			return nil, err
		}
	case "VlLs":
		size, err := reader.ReadInt()
		if err != nil {
			return nil, err
		}
		list := make([]*Item, 0, size)
		for i := 0; i < size; i++ {
			data, err := parseItem(reader)
			if err != nil {
				return nil, err
			}
			list = append(list, data)
		}
		item.Value = list
	case "doub":
		f, err := reader.ReadFloat64()
		if err != nil {
			return nil, err
		}
		item.Value = Double(f)
	case "UntF":
		uf := UnitFloat{}
		uf.Type, err = reader.ReadString(4)
		if err != nil {
			return nil, err
		}
		uf.Value, err = reader.ReadFloat64()
		if err != nil {
			return nil, err
		}
		item.Value = uf
	case "TEXT":
		str, err := reader.ReadUnicodeString()
		if err != nil {
			return nil, err
		}
		item.Value = Text(str)
	case "enum":
		enum := Enumerated{}
		enum.Type, err = reader.ReadDynamicString()
		if err != nil {
			return nil, err
		}
		enum.Value, err = reader.ReadDynamicString()
		if err != nil {
			return nil, err
		}
		item.Value = enum
	case "long":
		num, err := reader.ReadInt32()
		if err != nil {
			return nil, err
		}
		item.Value = Integer(num)
	case "comp":
		num, err := reader.ReadInt64()
		if err != nil {
			return nil, err
		}
		item.Value = LargeInteger(num)
	case "bool":
		b, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		item.Value = Boolean(b > 0)
	case "type", "GlbC":
		class := Class{}
		class.Name, err = reader.ReadUnicodeString()
		if err != nil {
			return nil, err
		}
		class.ID, err = reader.ReadDynamicString()
		if err != nil {
			return nil, err
		}
		item.Value = class
	case "alis":
		size, err := reader.ReadInt()
		if err != nil {
			return nil, err
		}
		str, err := reader.ReadString(size)
		if err != nil {
			return nil, err
		}
		item.Value = Alias(str)
	case "tdta":
		size, err := reader.ReadInt()
		if err != nil {
			return nil, err
		}
		buf, err := reader.ReadBytes(size)
		if err != nil {
			return nil, err
		}
		data, err := enginedata.Parser(buf[:])
		if err != nil {
			return nil, err
		}
		item.Value = data
	}

	return item, nil
}
