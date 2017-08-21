package descriptor

import (
	"fmt"
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

	Double float64

	UnitFloat struct {
		Type  string
		Value float64
	}

	Text string

	Enumerated struct {
		Type  string
		Value string
	}

	Integer int

	LargeInteger int64

	Boolean bool

	Class struct {
		Name string
		ID   string
	}

	Alias string

	Property struct {
		Name string
		ID   string
		Key  string
	}

	ReferenceEnum struct {
		Name string
		ID   string
		Type string
		Enum string
	}

	Offset struct {
		Name  string
		ID    string
		Value int32
	}
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

func Parse(reader *util.Reader) (*Descriptor, error) {
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
		item, err := parseItem(reader, true)
		if err != nil {
			return nil, err
		}
		items[item.Key] = item
	}
	obj.Items = items

	return obj, nil
}

func parseItem(reader *util.Reader, isKey bool) (*Item, error) {
	var err error

	item := &Item{}
	if isKey {
		item.Key, err = reader.ReadDynamicString()
		if err != nil {
			return nil, err
		}
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
		items := make([]Item, size)
		for i := range items {
			if items[i].Type, err = reader.ReadString(4); err != nil {
				return nil, err
			}
			switch items[i].Type {
			case "prop":
				name, err := reader.ReadUnicodeString()
				if err != nil {
					return nil, err
				}
				id, err := reader.ReadDynamicString()
				if err != nil {
					return nil, err
				}
				key, err := reader.ReadDynamicString()
				if err != nil {
					return nil, err
				}
				items[i].Value = Property{Name: name, ID: id, Key: key}
			case "Clss":
				class, err := parseClass(reader)
				if err != nil {
					return nil, err
				}
				items[i].Value = class
			case "Enmr":
				name, err := reader.ReadUnicodeString()
				if err != nil {
					return nil, err
				}
				id, err := reader.ReadDynamicString()
				if err != nil {
					return nil, err
				}
				typ, err := reader.ReadDynamicString()
				if err != nil {
					return nil, err
				}
				enum, err := reader.ReadDynamicString()
				if err != nil {
					return nil, err
				}
				items[i].Value = ReferenceEnum{Name: name, ID: id, Type: typ, Enum: enum}
			case "rele":
				name, err := reader.ReadUnicodeString()
				if err != nil {
					return nil, err
				}
				id, err := reader.ReadDynamicString()
				if err != nil {
					return nil, err
				}
				value, err := reader.ReadInt32()
				if err != nil {
					return nil, err
				}
				items[i].Value = Offset{Name: name, ID: id, Value: value}
			case "Idnt", "indx":
				id, err := reader.ReadInt()
				if err != nil {
					return nil, err
				}
				items[i].Value = id
			case "name":
				name, err := reader.ReadUnicodeString()
				if err != nil {
					return nil, err
				}
				items[i].Value = name
			}
		}
	case "Objc", "Glb0":
		item.Value, err = Parse(reader)
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
			data, err := parseItem(reader, false)
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
		class, err := parseClass(reader)
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
	default:
		panic(fmt.Sprintf("Unknown OSType key [%s] in entity [%s]", item.Key, item.Type))
	}

	return item, nil
}

func parseClass(reader *util.Reader) (*Class, error) {
	var err error
	class := &Class{}
	class.Name, err = reader.ReadUnicodeString()
	if err != nil {
		return nil, err
	}
	class.ID, err = reader.ReadDynamicString()
	if err != nil {
		return nil, err
	}
	return class, nil
}
