package flam

import (
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
)

type Bag map[string]any

func (bag *Bag) Clone() Bag {
	var cloner func(value any) any
	cloner = func(value any) any {
		switch typedValue := value.(type) {
		case []any:
			var result []any
			for _, i := range typedValue {
				result = append(result, cloner(i))
			}
			return result
		case Bag:
			return typedValue.Clone()
		case *Bag:
			return typedValue.Clone()
		default:
			return value
		}
	}

	target := Bag{}
	for key, value := range *bag {
		target[key] = cloner(value)
	}

	return target
}

func (bag *Bag) Entries() []string {
	var result []string
	for key := range *bag {
		result = append(result, key)
	}
	return result
}

func (bag *Bag) Has(
	path string,
) bool {
	_, e := bag.path(path)

	return e == nil
}

func (bag *Bag) Get(
	path string,
	def ...any,
) any {
	val, e := bag.path(path)
	if e != nil {
		if len(def) != 0 {
			return def[0]
		}

		return nil
	}

	return val
}

func (bag *Bag) Bool(
	path string,
	def ...bool,
) bool {
	switch val := bag.Get(path).(type) {
	case bool:
		return val
	default:
		if len(def) != 0 {
			return def[0]
		}
	}

	return false
}

func (bag *Bag) Int(
	path string,
	def ...int,
) int {
	switch val := bag.Get(path).(type) {
	case int:
		return val
	default:
		if len(def) != 0 {
			return def[0]
		}
	}

	return 0
}

func (bag *Bag) Int8(
	path string,
	def ...int8,
) int8 {
	switch val := bag.Get(path).(type) {
	case int8:
		return val
	default:
		if len(def) != 0 {
			return def[0]
		}
	}

	return 0
}

func (bag *Bag) Int16(
	path string,
	def ...int16,
) int16 {
	switch val := bag.Get(path).(type) {
	case int16:
		return val
	default:
		if len(def) != 0 {
			return def[0]
		}
	}

	return 0
}

func (bag *Bag) Int32(
	path string,
	def ...int32,
) int32 {
	switch val := bag.Get(path).(type) {
	case int32:
		return val
	default:
		if len(def) != 0 {
			return def[0]
		}
	}

	return 0
}

func (bag *Bag) Int64(
	path string,
	def ...int64,
) int64 {
	switch val := bag.Get(path).(type) {
	case int64:
		return val
	default:
		if len(def) != 0 {
			return def[0]
		}
	}

	return 0
}

func (bag *Bag) Uint(
	path string,
	def ...uint,
) uint {
	switch val := bag.Get(path).(type) {
	case uint:
		return val
	default:
		if len(def) != 0 {
			return def[0]
		}
	}

	return 0
}

func (bag *Bag) Uint8(
	path string,
	def ...uint8,
) uint8 {
	switch val := bag.Get(path).(type) {
	case uint8:
		return val
	default:
		if len(def) != 0 {
			return def[0]
		}
	}

	return 0
}

func (bag *Bag) Uint16(
	path string,
	def ...uint16,
) uint16 {
	switch val := bag.Get(path).(type) {
	case uint16:
		return val
	default:
		if len(def) != 0 {
			return def[0]
		}
	}

	return 0
}

func (bag *Bag) Uint32(
	path string,
	def ...uint32,
) uint32 {
	switch val := bag.Get(path).(type) {
	case uint32:
		return val
	default:
		if len(def) != 0 {
			return def[0]
		}
	}

	return 0
}

func (bag *Bag) Uint64(
	path string,
	def ...uint64,
) uint64 {
	switch val := bag.Get(path).(type) {
	case uint64:
		return val
	default:
		if len(def) != 0 {
			return def[0]
		}
	}

	return 0
}

func (bag *Bag) Float32(
	path string,
	def ...float32,
) float32 {
	switch val := bag.Get(path).(type) {
	case float32:
		return val
	default:
		if len(def) != 0 {
			return def[0]
		}
	}

	return 0
}

func (bag *Bag) Float64(
	path string,
	def ...float64,
) float64 {
	switch val := bag.Get(path).(type) {
	case float64:
		return val
	default:
		if len(def) != 0 {
			return def[0]
		}
	}

	return 0
}

func (bag *Bag) String(
	path string,
	def ...string,
) string {
	switch val := bag.Get(path).(type) {
	case string:
		return val
	default:
		if len(def) != 0 {
			return def[0]
		}
	}

	return ""
}

func (bag *Bag) StringMap(
	path string,
	def ...map[string]any,
) map[string]any {
	switch val := bag.Get(path).(type) {
	case map[string]any:
		return val
	default:
		if len(def) != 0 {
			return def[0]
		}
	}

	return nil
}

func (bag *Bag) StringMapString(
	path string,
	def ...map[string]string,
) map[string]string {
	switch val := bag.Get(path).(type) {
	case map[string]string:
		return val
	default:
		if len(def) != 0 {
			return def[0]
		}
	}

	return nil
}

func (bag *Bag) Slice(
	path string,
	def ...[]any,
) []any {
	switch val := bag.Get(path).(type) {
	case []any:
		return val
	default:
		if len(def) != 0 {
			return def[0]
		}
	}

	return nil
}

func (bag *Bag) StringSlice(
	path string,
	def ...[]string,
) []string {
	switch val := bag.Get(path).(type) {
	case []string:
		return val
	default:
		if len(def) != 0 {
			return def[0]
		}
	}

	return nil
}

func (bag *Bag) Duration(
	path string,
	def ...time.Duration,
) time.Duration {
	switch tval := bag.Get(path).(type) {
	case int:
		return time.Duration(tval) * time.Millisecond
	case int64:
		return time.Duration(tval) * time.Millisecond
	case time.Duration:
		return tval
	default:
		if len(def) != 0 {
			return def[0]
		}
	}

	return time.Duration(0)
}

func (bag *Bag) Bag(
	path string,
	def ...Bag,
) Bag {
	switch val := bag.Get(path).(type) {
	case Bag:
		return val
	default:
		if len(def) != 0 {
			return def[0]
		}
	}

	return nil
}

func (bag *Bag) Set(
	path string,
	value any,
) error {
	if path == "" {
		return newErrBagInvalidPath("")
	}

	parts := strings.Split(path, ".")
	it := bag
	if len(parts) == 1 {
		(*it)[path] = value
		return nil
	}

	generate := func(part string) {
		generate := false
		if next, ok := (*it)[part]; !ok {
			generate = true
		} else if _, ok = next.(Bag); !ok {
			generate = true
		}
		if generate {
			(*it)[part] = Bag{}
		}
	}

	for _, part := range parts[:len(parts)-1] {
		if part == "" {
			continue
		}

		generate(part)
		next := (*it)[part].(Bag)
		it = &next
	}

	part := parts[len(parts)-1:][0]
	generate(part)
	(*it)[part] = value

	return nil
}

func (bag *Bag) Merge(
	src Bag,
) *Bag {
	for key, value := range src {
		switch tValue := value.(type) {
		case Bag:
			switch tLocal := (*bag)[key].(type) {
			case Bag:
				tLocal.Merge(tValue)
			case *Bag:
				tLocal.Merge(tValue)
			default:
				v := Bag{}
				v.Merge(tValue)
				(*bag)[key] = v
			}
		case *Bag:
			switch tLocal := (*bag)[key].(type) {
			case Bag:
				tLocal.Merge(*tValue)
			case *Bag:
				tLocal.Merge(*tValue)
			default:
				v := Bag{}
				v.Merge(*tValue)
				(*bag)[key] = v
			}
		default:
			(*bag)[key] = value
		}
	}

	return bag
}

func (bag *Bag) Populate(
	target any,
	path ...string,
) error {
	p := ""
	if len(path) > 0 {
		p = path[0]
	}

	source := bag.Get(p, nil)
	if source == nil {
		return newErrBagInvalidPath(p)
	}

	return mapstructure.Decode(source, target)
}

func (bag *Bag) path(
	path string,
) (any, error) {
	var ok bool
	var it any

	it = *bag
	for _, part := range strings.Split(path, ".") {
		if part == "" {
			continue
		}

		switch typedIt := it.(type) {
		case Bag:
			if it, ok = typedIt[part]; !ok {
				return nil, newErrBagInvalidPath(path)
			}
		case *Bag:
			if it, ok = (*typedIt)[part]; !ok {
				return nil, newErrBagInvalidPath(path)
			}
		default:
			return nil, newErrBagInvalidPath(path)
		}
	}

	return it, nil
}
