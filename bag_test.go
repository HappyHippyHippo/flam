package flam

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Bag_Clone(t *testing.T) {
	scenarios := []struct {
		test     string
		bag      Bag
		expected Bag
	}{
		{
			test:     "should clone an empty bag",
			bag:      Bag{},
			expected: Bag{},
		},
		{
			test:     "should clone a simple bag",
			bag:      Bag{"field": 123},
			expected: Bag{"field": 123},
		},
		{
			test:     "should clone a nested bag",
			bag:      Bag{"field": Bag{"subfield": 456}},
			expected: Bag{"field": Bag{"subfield": 456}},
		},
		{
			test:     "should clone a nested bag (references)",
			bag:      Bag{"field": &Bag{"subfield": 456}},
			expected: Bag{"field": Bag{"subfield": 456}},
		},
		{
			test:     "should clone a bag with an array",
			bag:      Bag{"field": []any{1, 2, 3}},
			expected: Bag{"field": []any{1, 2, 3}},
		},
		{
			test:     "should clone a bag with a nested array",
			bag:      Bag{"field": []any{1, Bag{"subfield": 2}, 3}},
			expected: Bag{"field": []any{1, Bag{"subfield": 2}, 3}},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.test, func(t *testing.T) {
			result := scenario.bag.Clone()
			require.NotNil(t, result)
			require.NotSame(t, &result, &scenario.bag)
			assert.Equal(t, scenario.expected, result)
		})
	}
}

func Test_Bag_Entries(t *testing.T) {
	scenarios := []struct {
		test     string
		bag      Bag
		expected []string
	}{
		{
			test:     "should return nil for an empty bag",
			bag:      Bag{},
			expected: nil,
		},
		{
			test:     "should return a list of entries for a bag with one entry",
			bag:      Bag{"field": 123},
			expected: []string{"field"},
		},
		{
			test:     "should return a list of entries for a bag with multiple entries",
			bag:      Bag{"field1": 123, "field2": 456},
			expected: []string{"field1", "field2"},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.test, func(t *testing.T) {
			assert.ElementsMatch(t, scenario.expected, scenario.bag.Entries())
		})
	}
}

func Test_Bag_Has(t *testing.T) {
	scenarios := []struct {
		test     string
		bag      Bag
		path     string
		expected bool
	}{
		{
			test:     "should return true for an empty path",
			bag:      Bag{},
			path:     "",
			expected: true,
		},
		{
			test:     "should return false when checking for a path in an empty bag",
			bag:      Bag{},
			path:     "field",
			expected: false,
		},
		{
			test:     "should return true for a matching path",
			bag:      Bag{"field": 123},
			path:     "field",
			expected: true,
		},
		{
			test:     "should return false for a non-matching path",
			bag:      Bag{"field": 123},
			path:     "nonexistent",
			expected: false,
		},
		{
			test:     "should return true for a matching nested path",
			bag:      Bag{"field": Bag{"subfield": 456}},
			path:     "field.subfield",
			expected: true,
		},
		{
			test:     "should return false for a non-matching nested path",
			bag:      Bag{"field": Bag{"subfield": 456}},
			path:     "field.nonexistent",
			expected: false,
		},
		{
			test:     "should return true for a path with consecutive dots",
			bag:      Bag{"a": Bag{"b": 1}},
			path:     "a..b",
			expected: true,
		},
		{
			test:     "should return true for a path with leading dots",
			bag:      Bag{"a": 1},
			path:     ".a",
			expected: true,
		},
		{
			test:     "should return true for a path with trailing dots",
			bag:      Bag{"a": 1},
			path:     "a.",
			expected: true,
		},
		{
			test:     "should return false for a path through a non-bag value",
			bag:      Bag{"a": 1},
			path:     "a.b",
			expected: false,
		},
		{
			test:     "should return true for a path through a pointer to a bag",
			bag:      Bag{"a": &Bag{"b": 1}},
			path:     "a.b",
			expected: true,
		},
		{
			test:     "should return false for a non-matching path through a pointer to a bag",
			bag:      Bag{"a": &Bag{"b": 1}},
			path:     "a.c",
			expected: false,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.test, func(t *testing.T) {
			assert.Equal(t, scenario.expected, scenario.bag.Has(scenario.path))
		})
	}
}

func Test_Bag_Get(t *testing.T) {
	scenarios := []struct {
		test     string
		bag      Bag
		path     string
		def      []any
		expected any
	}{
		{
			test:     "should return an empty bag for an empty path on an empty bag",
			bag:      Bag{},
			path:     "",
			def:      nil,
			expected: Bag{},
		},
		{
			test:     "should return the bag itself for an empty path",
			bag:      Bag{"field": 123},
			path:     "",
			def:      nil,
			expected: Bag{"field": 123},
		},
		{
			test:     "should return a value for a valid path",
			bag:      Bag{"field": 123},
			path:     "field",
			def:      nil,
			expected: 123,
		},
		{
			test:     "should return a nested value for a valid path",
			bag:      Bag{"field": Bag{"subfield": 456}},
			path:     "field.subfield",
			def:      nil,
			expected: 456,
		},
		{
			test:     "should return nil for an invalid path without a default value",
			bag:      Bag{"field": 123},
			path:     "nonexistent",
			def:      nil,
			expected: nil,
		},
		{
			test:     "should return the bag itself for an empty path even with a default value",
			bag:      Bag{"field": 123},
			path:     "",
			def:      []any{"default"},
			expected: Bag{"field": 123},
		},
		{
			test:     "should return a value for a valid path even with a default value",
			bag:      Bag{"field": 123},
			path:     "field",
			def:      []any{"default"},
			expected: 123,
		},
		{
			test:     "should return the default value for an invalid path",
			bag:      Bag{"field": 123},
			path:     "nonexistent",
			def:      []any{"default"},
			expected: "default",
		},
		{
			test:     "should return the default value for a deep invalid path",
			bag:      Bag{"field": Bag{"subfield": 456}},
			path:     "field.nonexistent",
			def:      []any{"default"},
			expected: "default",
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.test, func(t *testing.T) {
			assert.Equal(t, scenario.expected, scenario.bag.Get(scenario.path, scenario.def...))
		})
	}
}

func Test_Bag_Bool(t *testing.T) {
	scenarios := []struct {
		test     string
		bag      Bag
		path     string
		def      []bool
		expected bool
	}{
		{
			test:     "should return false for an empty path without a default value",
			bag:      Bag{},
			path:     "",
			def:      nil,
			expected: false,
		},
		{
			test:     "should return the bool value for a valid path",
			bag:      Bag{"field": true},
			path:     "field",
			def:      nil,
			expected: true,
		},
		{
			test:     "should return false for a valid path with a non-bool value",
			bag:      Bag{"field": "value"},
			path:     "field",
			def:      nil,
			expected: false,
		},
		{
			test:     "should return the nested bool value for a valid path",
			bag:      Bag{"field": Bag{"subfield": true}},
			path:     "field.subfield",
			def:      nil,
			expected: true,
		},
		{
			test:     "should return false for an invalid path without a default value",
			bag:      Bag{"field": true},
			path:     "nonexistent",
			def:      nil,
			expected: false,
		},
		{
			test:     "should return the default value for an empty path",
			bag:      Bag{},
			path:     "",
			def:      []bool{true},
			expected: true,
		},
		{
			test:     "should return the bool value for a valid path when a default is provided",
			bag:      Bag{"field": true},
			path:     "field",
			def:      []bool{false},
			expected: true,
		},
		{
			test:     "should return the default value for a valid path with a non-bool value",
			bag:      Bag{"field": "value"},
			path:     "field",
			def:      []bool{true},
			expected: true,
		},
		{
			test:     "should return the default value for an invalid path",
			bag:      Bag{"field": true},
			path:     "nonexistent",
			def:      []bool{true},
			expected: true,
		},
		{
			test:     "should return the default value for a deep invalid path",
			bag:      Bag{"field": Bag{"subfield": true}},
			path:     "field.nonexistent",
			def:      []bool{true},
			expected: true,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.test, func(t *testing.T) {
			assert.Equal(t, scenario.expected, scenario.bag.Bool(scenario.path, scenario.def...))
		})
	}
}

func Test_Bag_Int(t *testing.T) {
	scenarios := []struct {
		test     string
		bag      Bag
		path     string
		def      []int
		expected int
	}{
		{
			test:     "should return 0 for an empty path without a default value",
			bag:      Bag{},
			path:     "",
			def:      nil,
			expected: 0,
		},
		{
			test:     "should return the int value for a valid path",
			bag:      Bag{"field": 123},
			path:     "field",
			def:      nil,
			expected: 123,
		},
		{
			test:     "should return 0 for a valid path with a non-int value",
			bag:      Bag{"field": "value"},
			path:     "field",
			def:      nil,
			expected: 0,
		},
		{
			test:     "should return the nested int value for a valid path",
			bag:      Bag{"field": Bag{"subfield": 123}},
			path:     "field.subfield",
			def:      nil,
			expected: 123,
		},
		{
			test:     "should return 0 for an invalid path without a default value",
			bag:      Bag{"field": 123},
			path:     "nonexistent",
			def:      nil,
			expected: 0,
		},
		{
			test:     "should return the default value for an empty path",
			bag:      Bag{},
			path:     "",
			def:      []int{456},
			expected: 456,
		},
		{
			test:     "should return the int value for a valid path when a default is provided",
			bag:      Bag{"field": 123},
			path:     "field",
			def:      []int{456},
			expected: 123,
		},
		{
			test:     "should return the default value for a valid path with a non-int value",
			bag:      Bag{"field": "value"},
			path:     "field",
			def:      []int{456},
			expected: 456,
		},
		{
			test:     "should return the default value for an invalid path",
			bag:      Bag{"field": 123},
			path:     "nonexistent",
			def:      []int{456},
			expected: 456,
		},
		{
			test:     "should return the default value for a deep invalid path",
			bag:      Bag{"field": Bag{"subfield": 123}},
			path:     "field.nonexistent",
			def:      []int{456},
			expected: 456,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.test, func(t *testing.T) {
			assert.Equal(t, scenario.expected, scenario.bag.Int(scenario.path, scenario.def...))
		})
	}
}

func Test_Bag_Int8(t *testing.T) {
	scenarios := []struct {
		test     string
		bag      Bag
		path     string
		def      []int8
		expected int8
	}{
		{
			test:     "should return 0 for an empty path without a default value",
			bag:      Bag{},
			path:     "",
			def:      nil,
			expected: 0,
		},
		{
			test:     "should return the int value for a valid path",
			bag:      Bag{"field": int8(123)},
			path:     "field",
			def:      nil,
			expected: int8(123),
		},
		{
			test:     "should return 0 for a valid path with a non-int value",
			bag:      Bag{"field": "value"},
			path:     "field",
			def:      nil,
			expected: 0,
		},
		{
			test:     "should return the nested int value for a valid path",
			bag:      Bag{"field": Bag{"subfield": int8(123)}},
			path:     "field.subfield",
			def:      nil,
			expected: int8(123),
		},
		{
			test:     "should return 0 for an invalid path without a default value",
			bag:      Bag{"field": 123},
			path:     "nonexistent",
			def:      nil,
			expected: 0,
		},
		{
			test:     "should return the default value for an empty path",
			bag:      Bag{},
			path:     "",
			def:      []int8{56},
			expected: 56,
		},
		{
			test:     "should return the int value for a valid path when a default is provided",
			bag:      Bag{"field": int8(123)},
			path:     "field",
			def:      []int8{56},
			expected: int8(123),
		},
		{
			test:     "should return the default value for a valid path with a non-int value",
			bag:      Bag{"field": "value"},
			path:     "field",
			def:      []int8{56},
			expected: 56,
		},
		{
			test:     "should return the default value for an invalid path",
			bag:      Bag{"field": 123},
			path:     "nonexistent",
			def:      []int8{56},
			expected: 56,
		},
		{
			test:     "should return the default value for a deep invalid path",
			bag:      Bag{"field": Bag{"subfield": 123}},
			path:     "field.nonexistent",
			def:      []int8{56},
			expected: 56,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.test, func(t *testing.T) {
			assert.Equal(t, scenario.expected, scenario.bag.Int8(scenario.path, scenario.def...))
		})
	}
}

func Test_Bag_Int16(t *testing.T) {
	scenarios := []struct {
		test     string
		bag      Bag
		path     string
		def      []int16
		expected int16
	}{
		{
			test:     "should return 0 for an empty path without a default value",
			bag:      Bag{},
			path:     "",
			def:      nil,
			expected: 0,
		},
		{
			test:     "should return the int value for a valid path",
			bag:      Bag{"field": int16(123)},
			path:     "field",
			def:      nil,
			expected: int16(123),
		},
		{
			test:     "should return 0 for a valid path with a non-int value",
			bag:      Bag{"field": "value"},
			path:     "field",
			def:      nil,
			expected: 0,
		},
		{
			test:     "should return the nested int value for a valid path",
			bag:      Bag{"field": Bag{"subfield": int16(123)}},
			path:     "field.subfield",
			def:      nil,
			expected: int16(123),
		},
		{
			test:     "should return 0 for an invalid path without a default value",
			bag:      Bag{"field": 123},
			path:     "nonexistent",
			def:      nil,
			expected: 0,
		},
		{
			test:     "should return the default value for an empty path",
			bag:      Bag{},
			path:     "",
			def:      []int16{56},
			expected: int16(56),
		},
		{
			test:     "should return the int value for a valid path when a default is provided",
			bag:      Bag{"field": int16(123)},
			path:     "field",
			def:      []int16{56},
			expected: int16(123),
		},
		{
			test:     "should return the default value for a valid path with a non-int value",
			bag:      Bag{"field": "value"},
			path:     "field",
			def:      []int16{56},
			expected: int16(56),
		},
		{
			test:     "should return the default value for an invalid path",
			bag:      Bag{"field": 123},
			path:     "nonexistent",
			def:      []int16{56},
			expected: int16(56),
		},
		{
			test:     "should return the default value for a deep invalid path",
			bag:      Bag{"field": Bag{"subfield": 123}},
			path:     "field.nonexistent",
			def:      []int16{56},
			expected: int16(56),
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.test, func(t *testing.T) {
			assert.Equal(t, scenario.expected, scenario.bag.Int16(scenario.path, scenario.def...))
		})
	}
}

func Test_Bag_Int32(t *testing.T) {
	scenarios := []struct {
		test     string
		bag      Bag
		path     string
		def      []int32
		expected int32
	}{
		{
			test:     "should return 0 for an empty path without a default value",
			bag:      Bag{},
			path:     "",
			def:      nil,
			expected: 0,
		},
		{
			test:     "should return the int value for a valid path",
			bag:      Bag{"field": int32(123)},
			path:     "field",
			def:      nil,
			expected: int32(123),
		},
		{
			test:     "should return 0 for a valid path with a non-int value",
			bag:      Bag{"field": "value"},
			path:     "field",
			def:      nil,
			expected: 0,
		},
		{
			test:     "should return the nested int value for a valid path",
			bag:      Bag{"field": Bag{"subfield": int32(123)}},
			path:     "field.subfield",
			def:      nil,
			expected: int32(123),
		},
		{
			test:     "should return 0 for an invalid path without a default value",
			bag:      Bag{"field": 123},
			path:     "nonexistent",
			def:      nil,
			expected: 0,
		},
		{
			test:     "should return the default value for an empty path",
			bag:      Bag{},
			path:     "",
			def:      []int32{56},
			expected: int32(56),
		},
		{
			test:     "should return the int value for a valid path when a default is provided",
			bag:      Bag{"field": int32(123)},
			path:     "field",
			def:      []int32{56},
			expected: int32(123),
		},
		{
			test:     "should return the default value for a valid path with a non-int value",
			bag:      Bag{"field": "value"},
			path:     "field",
			def:      []int32{56},
			expected: int32(56),
		},
		{
			test:     "should return the default value for an invalid path",
			bag:      Bag{"field": 123},
			path:     "nonexistent",
			def:      []int32{56},
			expected: int32(56),
		},
		{
			test:     "should return the default value for a deep invalid path",
			bag:      Bag{"field": Bag{"subfield": 123}},
			path:     "field.nonexistent",
			def:      []int32{56},
			expected: int32(56),
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.test, func(t *testing.T) {
			assert.Equal(t, scenario.expected, scenario.bag.Int32(scenario.path, scenario.def...))
		})
	}
}

func Test_Bag_Int64(t *testing.T) {
	scenarios := []struct {
		test     string
		bag      Bag
		path     string
		def      []int64
		expected int64
	}{
		{
			test:     "should return 0 for an empty path without a default value",
			bag:      Bag{},
			path:     "",
			def:      nil,
			expected: 0,
		},
		{
			test:     "should return the int value for a valid path",
			bag:      Bag{"field": int64(123)},
			path:     "field",
			def:      nil,
			expected: int64(123),
		},
		{
			test:     "should return 0 for a valid path with a non-int value",
			bag:      Bag{"field": "value"},
			path:     "field",
			def:      nil,
			expected: 0,
		},
		{
			test:     "should return the nested int value for a valid path",
			bag:      Bag{"field": Bag{"subfield": int64(123)}},
			path:     "field.subfield",
			def:      nil,
			expected: int64(123),
		},
		{
			test:     "should return 0 for an invalid path without a default value",
			bag:      Bag{"field": 123},
			path:     "nonexistent",
			def:      nil,
			expected: 0,
		},
		{
			test:     "should return the default value for an empty path",
			bag:      Bag{},
			path:     "",
			def:      []int64{56},
			expected: int64(56),
		},
		{
			test:     "should return the int value for a valid path when a default is provided",
			bag:      Bag{"field": int64(123)},
			path:     "field",
			def:      []int64{56},
			expected: int64(123),
		},
		{
			test:     "should return the default value for a valid path with a non-int value",
			bag:      Bag{"field": "value"},
			path:     "field",
			def:      []int64{56},
			expected: int64(56),
		},
		{
			test:     "should return the default value for an invalid path",
			bag:      Bag{"field": 123},
			path:     "nonexistent",
			def:      []int64{56},
			expected: int64(56),
		},
		{
			test:     "should return the default value for a deep invalid path",
			bag:      Bag{"field": Bag{"subfield": 123}},
			path:     "field.nonexistent",
			def:      []int64{56},
			expected: int64(56),
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.test, func(t *testing.T) {
			assert.Equal(t, scenario.expected, scenario.bag.Int64(scenario.path, scenario.def...))
		})
	}
}

func Test_Bag_Uint(t *testing.T) {
	scenarios := []struct {
		test     string
		bag      Bag
		path     string
		def      []uint
		expected uint
	}{
		{
			test:     "should return 0 for an empty path without a default value",
			bag:      Bag{},
			path:     "",
			def:      nil,
			expected: 0,
		},
		{
			test:     "should return the int value for a valid path",
			bag:      Bag{"field": uint(123)},
			path:     "field",
			def:      nil,
			expected: uint(123),
		},
		{
			test:     "should return 0 for a valid path with a non-int value",
			bag:      Bag{"field": "value"},
			path:     "field",
			def:      nil,
			expected: 0,
		},
		{
			test:     "should return the nested int value for a valid path",
			bag:      Bag{"field": Bag{"subfield": uint(123)}},
			path:     "field.subfield",
			def:      nil,
			expected: uint(123),
		},
		{
			test:     "should return 0 for an invalid path without a default value",
			bag:      Bag{"field": 123},
			path:     "nonexistent",
			def:      nil,
			expected: 0,
		},
		{
			test:     "should return the default value for an empty path",
			bag:      Bag{},
			path:     "",
			def:      []uint{56},
			expected: uint(56),
		},
		{
			test:     "should return the int value for a valid path when a default is provided",
			bag:      Bag{"field": uint(123)},
			path:     "field",
			def:      []uint{56},
			expected: uint(123),
		},
		{
			test:     "should return the default value for a valid path with a non-int value",
			bag:      Bag{"field": "value"},
			path:     "field",
			def:      []uint{56},
			expected: uint(56),
		},
		{
			test:     "should return the default value for an invalid path",
			bag:      Bag{"field": 123},
			path:     "nonexistent",
			def:      []uint{56},
			expected: uint(56),
		},
		{
			test:     "should return the default value for a deep invalid path",
			bag:      Bag{"field": Bag{"subfield": uint(123)}},
			path:     "field.nonexistent",
			def:      []uint{56},
			expected: uint(56),
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.test, func(t *testing.T) {
			assert.Equal(t, scenario.expected, scenario.bag.Uint(scenario.path, scenario.def...))
		})
	}
}

func Test_Bag_Uint8(t *testing.T) {
	scenarios := []struct {
		test     string
		bag      Bag
		path     string
		def      []uint8
		expected uint8
	}{
		{
			test:     "should return 0 for an empty path without a default value",
			bag:      Bag{},
			path:     "",
			def:      nil,
			expected: 0,
		},
		{
			test:     "should return the int value for a valid path",
			bag:      Bag{"field": uint8(123)},
			path:     "field",
			def:      nil,
			expected: uint8(123),
		},
		{
			test:     "should return 0 for a valid path with a non-int value",
			bag:      Bag{"field": "value"},
			path:     "field",
			def:      nil,
			expected: 0,
		},
		{
			test:     "should return the nested int value for a valid path",
			bag:      Bag{"field": Bag{"subfield": uint8(123)}},
			path:     "field.subfield",
			def:      nil,
			expected: uint8(123),
		},
		{
			test:     "should return 0 for an invalid path without a default value",
			bag:      Bag{"field": 123},
			path:     "nonexistent",
			def:      nil,
			expected: 0,
		},
		{
			test:     "should return the default value for an empty path",
			bag:      Bag{},
			path:     "",
			def:      []uint8{56},
			expected: uint8(56),
		},
		{
			test:     "should return the int value for a valid path when a default is provided",
			bag:      Bag{"field": uint8(123)},
			path:     "field",
			def:      []uint8{56},
			expected: uint8(123),
		},
		{
			test:     "should return the default value for a valid path with a non-int value",
			bag:      Bag{"field": "value"},
			path:     "field",
			def:      []uint8{56},
			expected: uint8(56),
		},
		{
			test:     "should return the default value for an invalid path",
			bag:      Bag{"field": 123},
			path:     "nonexistent",
			def:      []uint8{56},
			expected: uint8(56),
		},
		{
			test:     "should return the default value for a deep invalid path",
			bag:      Bag{"field": Bag{"subfield": uint8(123)}},
			path:     "field.nonexistent",
			def:      []uint8{56},
			expected: uint8(56),
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.test, func(t *testing.T) {
			assert.Equal(t, scenario.expected, scenario.bag.Uint8(scenario.path, scenario.def...))
		})
	}
}

func Test_Bag_Uint16(t *testing.T) {
	scenarios := []struct {
		test     string
		bag      Bag
		path     string
		def      []uint16
		expected uint16
	}{
		{
			test:     "should return 0 for an empty path without a default value",
			bag:      Bag{},
			path:     "",
			def:      nil,
			expected: 0,
		},
		{
			test:     "should return the int value for a valid path",
			bag:      Bag{"field": uint16(123)},
			path:     "field",
			def:      nil,
			expected: uint16(123),
		},
		{
			test:     "should return 0 for a valid path with a non-int value",
			bag:      Bag{"field": "value"},
			path:     "field",
			def:      nil,
			expected: 0,
		},
		{
			test:     "should return the nested int value for a valid path",
			bag:      Bag{"field": Bag{"subfield": uint16(123)}},
			path:     "field.subfield",
			def:      nil,
			expected: uint16(123),
		},
		{
			test:     "should return 0 for an invalid path without a default value",
			bag:      Bag{"field": 123},
			path:     "nonexistent",
			def:      nil,
			expected: 0,
		},
		{
			test:     "should return the default value for an empty path",
			bag:      Bag{},
			path:     "",
			def:      []uint16{56},
			expected: uint16(56),
		},
		{
			test:     "should return the int value for a valid path when a default is provided",
			bag:      Bag{"field": uint16(123)},
			path:     "field",
			def:      []uint16{56},
			expected: uint16(123),
		},
		{
			test:     "should return the default value for a valid path with a non-int value",
			bag:      Bag{"field": "value"},
			path:     "field",
			def:      []uint16{56},
			expected: uint16(56),
		},
		{
			test:     "should return the default value for an invalid path",
			bag:      Bag{"field": 123},
			path:     "nonexistent",
			def:      []uint16{56},
			expected: uint16(56),
		},
		{
			test:     "should return the default value for a deep invalid path",
			bag:      Bag{"field": Bag{"subfield": uint16(123)}},
			path:     "field.nonexistent",
			def:      []uint16{56},
			expected: uint16(56),
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.test, func(t *testing.T) {
			assert.Equal(t, scenario.expected, scenario.bag.Uint16(scenario.path, scenario.def...))
		})
	}
}

func Test_Bag_Uint32(t *testing.T) {
	scenarios := []struct {
		test     string
		bag      Bag
		path     string
		def      []uint32
		expected uint32
	}{
		{
			test:     "should return 0 for an empty path without a default value",
			bag:      Bag{},
			path:     "",
			def:      nil,
			expected: 0,
		},
		{
			test:     "should return the int value for a valid path",
			bag:      Bag{"field": uint32(123)},
			path:     "field",
			def:      nil,
			expected: uint32(123),
		},
		{
			test:     "should return 0 for a valid path with a non-int value",
			bag:      Bag{"field": "value"},
			path:     "field",
			def:      nil,
			expected: 0,
		},
		{
			test:     "should return the nested int value for a valid path",
			bag:      Bag{"field": Bag{"subfield": uint32(123)}},
			path:     "field.subfield",
			def:      nil,
			expected: uint32(123),
		},
		{
			test:     "should return 0 for an invalid path without a default value",
			bag:      Bag{"field": 123},
			path:     "nonexistent",
			def:      nil,
			expected: 0,
		},
		{
			test:     "should return the default value for an empty path",
			bag:      Bag{},
			path:     "",
			def:      []uint32{56},
			expected: uint32(56),
		},
		{
			test:     "should return the int value for a valid path when a default is provided",
			bag:      Bag{"field": uint32(123)},
			path:     "field",
			def:      []uint32{56},
			expected: uint32(123),
		},
		{
			test:     "should return the default value for a valid path with a non-int value",
			bag:      Bag{"field": "value"},
			path:     "field",
			def:      []uint32{56},
			expected: uint32(56),
		},
		{
			test:     "should return the default value for an invalid path",
			bag:      Bag{"field": 123},
			path:     "nonexistent",
			def:      []uint32{56},
			expected: uint32(56),
		},
		{
			test:     "should return the default value for a deep invalid path",
			bag:      Bag{"field": Bag{"subfield": uint32(123)}},
			path:     "field.nonexistent",
			def:      []uint32{56},
			expected: uint32(56),
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.test, func(t *testing.T) {
			assert.Equal(t, scenario.expected, scenario.bag.Uint32(scenario.path, scenario.def...))
		})
	}
}

func Test_Bag_Uint64(t *testing.T) {
	scenarios := []struct {
		test     string
		bag      Bag
		path     string
		def      []uint64
		expected uint64
	}{
		{
			test:     "should return 0 for an empty path without a default value",
			bag:      Bag{},
			path:     "",
			def:      nil,
			expected: 0,
		},
		{
			test:     "should return the int value for a valid path",
			bag:      Bag{"field": uint64(123)},
			path:     "field",
			def:      nil,
			expected: uint64(123),
		},
		{
			test:     "should return 0 for a valid path with a non-int value",
			bag:      Bag{"field": "value"},
			path:     "field",
			def:      nil,
			expected: 0,
		},
		{
			test:     "should return the nested int value for a valid path",
			bag:      Bag{"field": Bag{"subfield": uint64(123)}},
			path:     "field.subfield",
			def:      nil,
			expected: uint64(123),
		},
		{
			test:     "should return 0 for an invalid path without a default value",
			bag:      Bag{"field": 123},
			path:     "nonexistent",
			def:      nil,
			expected: 0,
		},
		{
			test:     "should return the default value for an empty path",
			bag:      Bag{},
			path:     "",
			def:      []uint64{56},
			expected: uint64(56),
		},
		{
			test:     "should return the int value for a valid path when a default is provided",
			bag:      Bag{"field": uint64(123)},
			path:     "field",
			def:      []uint64{56},
			expected: uint64(123),
		},
		{
			test:     "should return the default value for a valid path with a non-int value",
			bag:      Bag{"field": "value"},
			path:     "field",
			def:      []uint64{56},
			expected: uint64(56),
		},
		{
			test:     "should return the default value for an invalid path",
			bag:      Bag{"field": 123},
			path:     "nonexistent",
			def:      []uint64{56},
			expected: uint64(56),
		},
		{
			test:     "should return the default value for a deep invalid path",
			bag:      Bag{"field": Bag{"subfield": uint64(123)}},
			path:     "field.nonexistent",
			def:      []uint64{56},
			expected: uint64(56),
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.test, func(t *testing.T) {
			assert.Equal(t, scenario.expected, scenario.bag.Uint64(scenario.path, scenario.def...))
		})
	}
}

func Test_Bag_Float32(t *testing.T) {
	scenarios := []struct {
		test     string
		bag      Bag
		path     string
		def      []float32
		expected float32
	}{
		{
			test:     "should return 0.0 for an empty path without a default value",
			bag:      Bag{},
			path:     "",
			def:      nil,
			expected: float32(0.0),
		},
		{
			test:     "should return the float value for a valid path",
			bag:      Bag{"field": float32(123.45)},
			path:     "field",
			def:      nil,
			expected: float32(123.45),
		},
		{
			test:     "should return 0.0 for a valid path with a non-float value",
			bag:      Bag{"field": "value"},
			path:     "field",
			def:      nil,
			expected: float32(0.0),
		},
		{
			test:     "should return the nested float value for a valid path",
			bag:      Bag{"field": Bag{"subfield": float32(123.45)}},
			path:     "field.subfield",
			def:      nil,
			expected: float32(123.45),
		},
		{
			test:     "should return 0.0 for an invalid path without a default value",
			bag:      Bag{"field": float32(123.45)},
			path:     "nonexistent",
			def:      nil,
			expected: float32(0.0),
		},
		{
			test:     "should return the default value for an empty path",
			bag:      Bag{},
			path:     "",
			def:      []float32{456.78},
			expected: float32(456.78),
		},
		{
			test:     "should return the float value for a valid path when a default is provided",
			bag:      Bag{"field": float32(123.45)},
			path:     "field",
			def:      []float32{456.78},
			expected: float32(123.45),
		},
		{
			test:     "should return the default value for a valid path with a non-float value",
			bag:      Bag{"field": "value"},
			path:     "field",
			def:      []float32{456.78},
			expected: float32(456.78),
		},
		{
			test:     "should return the default value for an invalid path",
			bag:      Bag{"field": float32(123.45)},
			path:     "nonexistent",
			def:      []float32{456.78},
			expected: float32(456.78),
		},
		{
			test:     "should return the default value for a deep invalid path",
			bag:      Bag{"field": Bag{"subfield": float32(123.45)}},
			path:     "field.nonexistent",
			def:      []float32{456.78},
			expected: float32(456.78),
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.test, func(t *testing.T) {
			assert.Equal(t, scenario.expected, scenario.bag.Float32(scenario.path, scenario.def...))
		})
	}
}

func Test_Bag_Float64(t *testing.T) {
	scenarios := []struct {
		test     string
		bag      Bag
		path     string
		def      []float64
		expected float64
	}{
		{
			test:     "should return 0.0 for an empty path without a default value",
			bag:      Bag{},
			path:     "",
			def:      nil,
			expected: 0.0,
		},
		{
			test:     "should return the float value for a valid path",
			bag:      Bag{"field": 123.45},
			path:     "field",
			def:      nil,
			expected: 123.45,
		},
		{
			test:     "should return 0.0 for a valid path with a non-float value",
			bag:      Bag{"field": "value"},
			path:     "field",
			def:      nil,
			expected: 0.0,
		},
		{
			test:     "should return the nested float value for a valid path",
			bag:      Bag{"field": Bag{"subfield": 123.45}},
			path:     "field.subfield",
			def:      nil,
			expected: 123.45,
		},
		{
			test:     "should return 0.0 for an invalid path without a default value",
			bag:      Bag{"field": 123.45},
			path:     "nonexistent",
			def:      nil,
			expected: 0.0,
		},
		{
			test:     "should return the default value for an empty path",
			bag:      Bag{},
			path:     "",
			def:      []float64{456.78},
			expected: 456.78,
		},
		{
			test:     "should return the float value for a valid path when a default is provided",
			bag:      Bag{"field": 123.45},
			path:     "field",
			def:      []float64{456.78},
			expected: 123.45,
		},
		{
			test:     "should return the default value for a valid path with a non-float value",
			bag:      Bag{"field": "value"},
			path:     "field",
			def:      []float64{456.78},
			expected: 456.78,
		},
		{
			test:     "should return the default value for an invalid path",
			bag:      Bag{"field": 123.45},
			path:     "nonexistent",
			def:      []float64{456.78},
			expected: 456.78,
		},
		{
			test:     "should return the default value for a deep invalid path",
			bag:      Bag{"field": Bag{"subfield": 123.45}},
			path:     "field.nonexistent",
			def:      []float64{456.78},
			expected: 456.78,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.test, func(t *testing.T) {
			assert.Equal(t, scenario.expected, scenario.bag.Float64(scenario.path, scenario.def...))
		})
	}
}

func Test_Bag_String(t *testing.T) {
	scenarios := []struct {
		test     string
		bag      Bag
		path     string
		def      []string
		expected string
	}{
		{
			test:     "should return empty string for an empty path without a default value",
			bag:      Bag{},
			path:     "",
			def:      nil,
			expected: "",
		},
		{
			test:     "should return the string value for a valid path",
			bag:      Bag{"field": "value"},
			path:     "field",
			def:      nil,
			expected: "value",
		},
		{
			test:     "should return empty string for a valid path with a non-string value",
			bag:      Bag{"field": 123},
			path:     "field",
			def:      nil,
			expected: "",
		},
		{
			test:     "should return the nested string value for a valid path",
			bag:      Bag{"field": Bag{"subfield": "value"}},
			path:     "field.subfield",
			def:      nil,
			expected: "value",
		},
		{
			test:     "should return empty string for an invalid path without a default value",
			bag:      Bag{"field": "value"},
			path:     "nonexistent",
			def:      nil,
			expected: "",
		},
		{
			test:     "should return the default value for an empty path",
			bag:      Bag{},
			path:     "",
			def:      []string{"default"},
			expected: "default",
		},
		{
			test:     "should return the string value for a valid path when a default is provided",
			bag:      Bag{"field": "value"},
			path:     "field",
			def:      []string{"default"},
			expected: "value",
		},
		{
			test:     "should return the default value for a valid path with a non-string value",
			bag:      Bag{"field": 123},
			path:     "field",
			def:      []string{"default"},
			expected: "default",
		},
		{
			test:     "should return the default value for an invalid path",
			bag:      Bag{"field": "value"},
			path:     "nonexistent",
			def:      []string{"default"},
			expected: "default",
		},
		{
			test:     "should return the default value for a deep invalid path",
			bag:      Bag{"field": Bag{"subfield": "value"}},
			path:     "field.nonexistent",
			def:      []string{"default"},
			expected: "default",
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.test, func(t *testing.T) {
			assert.Equal(t, scenario.expected, scenario.bag.String(scenario.path, scenario.def...))
		})
	}
}

func Test_Bag_StringMap(t *testing.T) {
	scenarios := []struct {
		test     string
		bag      Bag
		path     string
		def      []map[string]any
		expected map[string]any
	}{
		{
			test:     "should return nil for an empty path without a default value",
			bag:      Bag{},
			path:     "",
			def:      nil,
			expected: nil,
		},
		{
			test:     "should return the string map value for a valid path",
			bag:      Bag{"field": map[string]any{"a": 1}},
			path:     "field",
			def:      nil,
			expected: map[string]any{"a": 1},
		},
		{
			test:     "should return nil for a valid path with a non-string value",
			bag:      Bag{"field": 123},
			path:     "field",
			def:      nil,
			expected: nil,
		},
		{
			test:     "should return the nested string map value for a valid path",
			bag:      Bag{"field": Bag{"subfield": map[string]any{"a": 1}}},
			path:     "field.subfield",
			def:      nil,
			expected: map[string]any{"a": 1},
		},
		{
			test:     "should return nil for an invalid path without a default value",
			bag:      Bag{"field": "value"},
			path:     "nonexistent",
			def:      nil,
			expected: nil,
		},
		{
			test:     "should return the default value for an empty path",
			bag:      Bag{},
			path:     "",
			def:      []map[string]any{{"default": 1}},
			expected: map[string]any{"default": 1},
		},
		{
			test:     "should return the string map value for a valid path when a default is provided",
			bag:      Bag{"field": map[string]any{"a": 1}},
			path:     "field",
			def:      []map[string]any{{"default": 1}},
			expected: map[string]any{"a": 1},
		},
		{
			test:     "should return the default value for a valid path with a non-string map value",
			bag:      Bag{"field": 123},
			path:     "field",
			def:      []map[string]any{{"default": 1}},
			expected: map[string]any{"default": 1},
		},
		{
			test:     "should return the default value for an invalid path",
			bag:      Bag{"field": "value"},
			path:     "nonexistent",
			def:      []map[string]any{{"default": 1}},
			expected: map[string]any{"default": 1},
		},
		{
			test:     "should return the default value for a deep invalid path",
			bag:      Bag{"field": Bag{"subfield": "value"}},
			path:     "field.nonexistent",
			def:      []map[string]any{{"default": 1}},
			expected: map[string]any{"default": 1},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.test, func(t *testing.T) {
			assert.Equal(t, scenario.expected, scenario.bag.StringMap(scenario.path, scenario.def...))
		})
	}
}

func Test_Bag_StringMapString(t *testing.T) {
	scenarios := []struct {
		test     string
		bag      Bag
		path     string
		def      []map[string]string
		expected map[string]string
	}{
		{
			test:     "should return nil for an empty path without a default value",
			bag:      Bag{},
			path:     "",
			def:      nil,
			expected: nil,
		},
		{
			test:     "should return the string map value for a valid path",
			bag:      Bag{"field": map[string]string{"a": "1"}},
			path:     "field",
			def:      nil,
			expected: map[string]string{"a": "1"},
		},
		{
			test:     "should return nil for a valid path with a non-string value",
			bag:      Bag{"field": 123},
			path:     "field",
			def:      nil,
			expected: nil,
		},
		{
			test:     "should return the nested string map value for a valid path",
			bag:      Bag{"field": Bag{"subfield": map[string]string{"a": "1"}}},
			path:     "field.subfield",
			def:      nil,
			expected: map[string]string{"a": "1"},
		},
		{
			test:     "should return nil for an invalid path without a default value",
			bag:      Bag{"field": "value"},
			path:     "nonexistent",
			def:      nil,
			expected: nil,
		},
		{
			test:     "should return the default value for an empty path",
			bag:      Bag{},
			path:     "",
			def:      []map[string]string{{"default": "1"}},
			expected: map[string]string{"default": "1"},
		},
		{
			test:     "should return the string map value for a valid path when a default is provided",
			bag:      Bag{"field": map[string]string{"a": "1"}},
			path:     "field",
			def:      []map[string]string{{"default": "1"}},
			expected: map[string]string{"a": "1"},
		},
		{
			test:     "should return the default value for a valid path with a non-string map value",
			bag:      Bag{"field": 123},
			path:     "field",
			def:      []map[string]string{{"default": "1"}},
			expected: map[string]string{"default": "1"},
		},
		{
			test:     "should return the default value for an invalid path",
			bag:      Bag{"field": "value"},
			path:     "nonexistent",
			def:      []map[string]string{{"default": "1"}},
			expected: map[string]string{"default": "1"},
		},
		{
			test:     "should return the default value for a deep invalid path",
			bag:      Bag{"field": Bag{"subfield": "value"}},
			path:     "field.nonexistent",
			def:      []map[string]string{{"default": "1"}},
			expected: map[string]string{"default": "1"},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.test, func(t *testing.T) {
			assert.Equal(t, scenario.expected, scenario.bag.StringMapString(scenario.path, scenario.def...))
		})
	}
}

func Test_Bag_Slice(t *testing.T) {
	scenarios := []struct {
		test     string
		bag      Bag
		path     string
		def      [][]any
		expected []any
	}{
		{
			test:     "should return nil for an empty path without a default value",
			bag:      Bag{},
			path:     "",
			def:      nil,
			expected: nil,
		},
		{
			test:     "should return the slice value for a valid path",
			bag:      Bag{"field": []any{"a", 1}},
			path:     "field",
			def:      nil,
			expected: []any{"a", 1},
		},
		{
			test:     "should return nil for a valid path with a non-string value",
			bag:      Bag{"field": 123},
			path:     "field",
			def:      nil,
			expected: nil,
		},
		{
			test:     "should return the nested slice value for a valid path",
			bag:      Bag{"field": []any{"a", 1}},
			path:     "field",
			def:      nil,
			expected: []any{"a", 1},
		},
		{
			test:     "should return nil for an invalid path without a default value",
			bag:      Bag{"field": "value"},
			path:     "nonexistent",
			def:      nil,
			expected: nil,
		},
		{
			test:     "should return the default value for an empty path",
			bag:      Bag{},
			path:     "",
			def:      [][]any{{1}},
			expected: []any{1},
		},
		{
			test:     "should return the slice value for a valid path when a default is provided",
			bag:      Bag{"field": []any{"a", 1}},
			path:     "field",
			def:      [][]any{{1}},
			expected: []any{"a", 1},
		},
		{
			test:     "should return the default value for a valid path with a non-string map value",
			bag:      Bag{"field": 123},
			path:     "field",
			def:      [][]any{{1}},
			expected: []any{1},
		},
		{
			test:     "should return the default value for an invalid path",
			bag:      Bag{"field": "value"},
			path:     "nonexistent",
			def:      [][]any{{1}},
			expected: []any{1},
		},
		{
			test:     "should return the default value for a deep invalid path",
			bag:      Bag{"field": Bag{"subfield": "value"}},
			path:     "field.nonexistent",
			def:      [][]any{{1}},
			expected: []any{1},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.test, func(t *testing.T) {
			assert.Equal(t, scenario.expected, scenario.bag.Slice(scenario.path, scenario.def...))
		})
	}
}

func Test_Bag_StringSlice(t *testing.T) {
	scenarios := []struct {
		test     string
		bag      Bag
		path     string
		def      [][]string
		expected []string
	}{
		{
			test:     "should return nil for an empty path without a default value",
			bag:      Bag{},
			path:     "",
			def:      nil,
			expected: nil,
		},
		{
			test:     "should return the string slice value for a valid path",
			bag:      Bag{"field": []string{"a", "b"}},
			path:     "field",
			def:      nil,
			expected: []string{"a", "b"},
		},
		{
			test:     "should return nil for a valid path with a non-slice value",
			bag:      Bag{"field": "value"},
			path:     "field",
			def:      nil,
			expected: nil,
		},
		{
			test:     "should return nil for a valid path with a slice of non-string",
			bag:      Bag{"field": []int{1, 2}},
			path:     "field",
			def:      nil,
			expected: nil,
		},
		{
			test:     "should return the nested string slice value for a valid path",
			bag:      Bag{"field": Bag{"subfield": []string{"a", "b"}}},
			path:     "field.subfield",
			def:      nil,
			expected: []string{"a", "b"},
		},
		{
			test:     "should return nil for an invalid path without a default value",
			bag:      Bag{"field": []string{"a", "b"}},
			path:     "nonexistent",
			def:      nil,
			expected: nil,
		},
		{
			test:     "should return the default value for an empty path",
			bag:      Bag{},
			path:     "",
			def:      [][]string{{"c", "d"}},
			expected: []string{"c", "d"},
		},
		{
			test:     "should return the string slice value for a valid path when a default is provided",
			bag:      Bag{"field": []string{"a", "b"}},
			path:     "field",
			def:      [][]string{{"c", "d"}},
			expected: []string{"a", "b"},
		},
		{
			test:     "should return the default value for a valid path with a non-slice value",
			bag:      Bag{"field": "value"},
			path:     "field",
			def:      [][]string{{"c", "d"}},
			expected: []string{"c", "d"},
		},
		{
			test:     "should return the default value for an invalid path",
			bag:      Bag{"field": []string{"a", "b"}},
			path:     "nonexistent",
			def:      [][]string{{"c", "d"}},
			expected: []string{"c", "d"},
		},
		{
			test:     "should return the default value for a deep invalid path",
			bag:      Bag{"field": Bag{"subfield": []string{"a", "b"}}},
			path:     "field.nonexistent",
			def:      [][]string{{"c", "d"}},
			expected: []string{"c", "d"},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.test, func(t *testing.T) {
			assert.Equal(t, scenario.expected, scenario.bag.StringSlice(scenario.path, scenario.def...))
		})
	}
}

func Test_Bag_Duration(t *testing.T) {
	defaultValue := time.Minute

	scenarios := []struct {
		test     string
		bag      Bag
		path     string
		def      []time.Duration
		expected time.Duration
	}{
		{
			test:     "should return 0 for an empty path without a default value",
			bag:      Bag{},
			path:     "",
			def:      nil,
			expected: 0,
		},
		{
			test:     "should return the duration value for a valid path",
			bag:      Bag{"field": time.Second},
			path:     "field",
			def:      nil,
			expected: time.Second,
		},
		{
			test:     "should return 0 for a valid path with a non-duration value",
			bag:      Bag{"field": "value"},
			path:     "field",
			def:      nil,
			expected: 0,
		},
		{
			test:     "should return the duration value for a valid path",
			bag:      Bag{"field": time.Second},
			path:     "field",
			def:      nil,
			expected: time.Second,
		},
		{
			test:     "should return the duration value for a valid path (int convertion)",
			bag:      Bag{"field": 1000},
			path:     "field",
			def:      nil,
			expected: 1000 * time.Millisecond,
		},
		{
			test:     "should return the duration value for a valid path (int64 convertion)",
			bag:      Bag{"field": int64(1000)},
			path:     "field",
			def:      nil,
			expected: 1000 * time.Millisecond,
		},
		{
			test:     "should return 0 for an invalid path without a default value",
			bag:      Bag{"field": "value"},
			path:     "nonexistent",
			def:      nil,
			expected: 0,
		},
		{
			test:     "should return the default value for an empty path",
			bag:      Bag{},
			path:     "",
			def:      []time.Duration{defaultValue},
			expected: defaultValue,
		},
		{
			test:     "should return the duration value for a valid path when a default is provided",
			bag:      Bag{"field": time.Second},
			path:     "field",
			def:      []time.Duration{defaultValue},
			expected: time.Second,
		},
		{
			test:     "should return the default value for a valid path with a non-duration value",
			bag:      Bag{"field": "value"},
			path:     "field",
			def:      []time.Duration{defaultValue},
			expected: defaultValue,
		},
		{
			test:     "should return the default value for an invalid path",
			bag:      Bag{"field": "value"},
			path:     "nonexistent",
			def:      []time.Duration{defaultValue},
			expected: defaultValue,
		},
		{
			test:     "should return the default value for a deep invalid path",
			bag:      Bag{"field": Bag{"subfield": "value"}},
			path:     "field.nonexistent",
			def:      []time.Duration{defaultValue},
			expected: defaultValue,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.test, func(t *testing.T) {
			assert.Equal(t, scenario.expected, scenario.bag.Duration(scenario.path, scenario.def...))
		})
	}
}

func Test_Bag_Bag(t *testing.T) {
	scenarios := []struct {
		test     string
		bag      Bag
		path     string
		def      []Bag
		expected Bag
	}{
		{
			test:     "should return a copy for an empty path without a default value",
			bag:      Bag{"field": 123},
			path:     "",
			def:      nil,
			expected: Bag{"field": 123},
		},
		{
			test:     "should return the bag value for a valid path",
			bag:      Bag{"field": Bag{"subfield": "value"}},
			path:     "field",
			def:      nil,
			expected: Bag{"subfield": "value"},
		},
		{
			test:     "should return nil for a valid path with a non-bag value",
			bag:      Bag{"field": 123},
			path:     "field",
			def:      nil,
			expected: nil,
		},
		{
			test:     "should return nil for an invalid path without a default value",
			bag:      Bag{"field": "value"},
			path:     "nonexistent",
			def:      nil,
			expected: nil,
		},
		{
			test:     "should return the bag value for a valid path when a default is provided",
			bag:      Bag{"field": Bag{"subfield": "value"}},
			path:     "field",
			def:      []Bag{{"default": "value"}},
			expected: Bag{"subfield": "value"},
		},
		{
			test:     "should return the default value for a valid path with a non-bag value",
			bag:      Bag{"field": 123},
			path:     "field",
			def:      []Bag{{"default": "value"}},
			expected: Bag{"default": "value"},
		},
		{
			test:     "should return the default value for an invalid path",
			bag:      Bag{"field": "value"},
			path:     "nonexistent",
			def:      []Bag{{"default": "value"}},
			expected: Bag{"default": "value"},
		},
		{
			test:     "should return the default value for a deep invalid path",
			bag:      Bag{"field": Bag{"subfield": "value"}},
			path:     "field.nonexistent",
			def:      []Bag{{"default": "value"}},
			expected: Bag{"default": "value"},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.test, func(t *testing.T) {
			assert.Equal(t, scenario.expected, scenario.bag.Bag(scenario.path, scenario.def...))
		})
	}
}

func Test_Bag_Set(t *testing.T) {
	scenarios := []struct {
		test        string
		bag         Bag
		path        string
		value       any
		expected    Bag
		expectedErr error
	}{
		{
			test:     "should set a value at the top level",
			bag:      Bag{},
			path:     "field",
			value:    123,
			expected: Bag{"field": 123},
		},
		{
			test:     "should set a value in a nested bag",
			bag:      Bag{"nested": Bag{}},
			path:     "nested.field",
			value:    "value",
			expected: Bag{"nested": Bag{"field": "value"}},
		},
		{
			test:     "should set a value in a nested bag (dot sequence)",
			bag:      Bag{"nested": Bag{}},
			path:     "nested...field",
			value:    "value",
			expected: Bag{"nested": Bag{"field": "value"}},
		},
		{
			test:     "should create nested bags to set a value",
			bag:      Bag{},
			path:     "a.b.c",
			value:    true,
			expected: Bag{"a": Bag{"b": Bag{"c": true}}},
		},
		{
			test:     "should overwrite an existing value",
			bag:      Bag{"field": 123},
			path:     "field",
			value:    "new_value",
			expected: Bag{"field": "new_value"},
		},
		{
			test:        "should return an error for an empty path",
			bag:         Bag{},
			path:        "",
			value:       123,
			expectedErr: ErrBagInvalidPath,
		},
		{
			test:     "should overwrite a non-bag value with a bag",
			bag:      Bag{"a": 123},
			path:     "a.b",
			value:    "hello",
			expected: Bag{"a": Bag{"b": "hello"}},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.test, func(t *testing.T) {
			e := scenario.bag.Set(scenario.path, scenario.value)

			if scenario.expectedErr != nil {
				assert.ErrorIs(t, e, scenario.expectedErr)
				return
			}

			assert.NoError(t, e)
			assert.Equal(t, scenario.expected, scenario.bag)
		})
	}
}

func Test_Bag_Merge(t *testing.T) {
	scenarios := []struct {
		test     string
		dest     Bag
		src      Bag
		expected Bag
	}{
		{
			test:     "should merge simple values into an empty bag",
			dest:     Bag{},
			src:      Bag{"a": 1, "b": "hello"},
			expected: Bag{"a": 1, "b": "hello"},
		},
		{
			test:     "should add new values without overwriting existing ones",
			dest:     Bag{"a": 1},
			src:      Bag{"b": 2},
			expected: Bag{"a": 1, "b": 2},
		},
		{
			test:     "should overwrite existing values",
			dest:     Bag{"a": 1, "b": "old"},
			src:      Bag{"b": "new", "c": 3},
			expected: Bag{"a": 1, "b": "new", "c": 3},
		},
		{
			test:     "should merge nested bags",
			dest:     Bag{"nested": Bag{"a": 1}},
			src:      Bag{"nested": Bag{"b": 2}},
			expected: Bag{"nested": Bag{"a": 1, "b": 2}},
		},
		{
			test:     "should merge nested bags (pointer in src)",
			dest:     Bag{"nested": Bag{"a": 1}},
			src:      Bag{"nested": &Bag{"b": 2}},
			expected: Bag{"nested": Bag{"a": 1, "b": 2}},
		},
		{
			test:     "should create nested bag if destination has a non-bag value",
			dest:     Bag{"a": 123},
			src:      Bag{"a": Bag{"b": "hello"}},
			expected: Bag{"a": Bag{"b": "hello"}},
		},
		{
			test:     "should handle complex nested merging",
			dest:     Bag{"a": Bag{"b": 1}, "c": "foo"},
			src:      Bag{"a": Bag{"d": 2}, "c": "bar", "e": Bag{"f": 3}},
			expected: Bag{"a": Bag{"b": 1, "d": 2}, "c": "bar", "e": Bag{"f": 3}},
		},
		{
			test:     "should merge into a nested pointer bag",
			dest:     Bag{"nested": &Bag{"a": 1}},
			src:      Bag{"nested": Bag{"b": 2}},
			expected: Bag{"nested": &Bag{"a": 1, "b": 2}},
		},
		{
			test:     "should merge pointer bag into a nested pointer bag",
			dest:     Bag{"nested": &Bag{"a": 1}},
			src:      Bag{"nested": &Bag{"b": 2}},
			expected: Bag{"nested": &Bag{"a": 1, "b": 2}},
		},
		{
			test:     "should create nested bag from pointer if destination has a non-bag value",
			dest:     Bag{"a": 123},
			src:      Bag{"a": &Bag{"b": "hello"}},
			expected: Bag{"a": Bag{"b": "hello"}},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.test, func(t *testing.T) {
			scenario.dest.Merge(scenario.src)
			assert.Equal(t, scenario.expected, scenario.dest)
		})
	}
}

func Test_Bag_Populate(t *testing.T) {
	type simpleStruct struct {
		Field int
	}

	type complexStruct struct {
		Name     string `mapstructure:"name"`
		Value    int    `mapstructure:"value"`
		Nested   simpleStruct
		Children []string
	}

	t.Run("without path", func(t *testing.T) {
		scenarios := []struct {
			test        string
			bag         Bag
			target      any
			expected    any
			expectedErr error
		}{
			{
				test:     "should populate a struct with a simple scalar value",
				bag:      Bag{"field": 123},
				target:   &simpleStruct{},
				expected: &simpleStruct{Field: 123},
			},
			{
				test: "should populate a complex struct with tags",
				bag: Bag{
					"name":  "test_name",
					"value": 999,
					"Nested": Bag{
						"Field": 789,
					},
					"Children": []any{"child1", "child2"},
				},
				target: &complexStruct{},
				expected: &complexStruct{
					Name:     "test_name",
					Value:    999,
					Nested:   simpleStruct{Field: 789},
					Children: []string{"child1", "child2"},
				},
			},
		}

		for _, scenario := range scenarios {
			t.Run(scenario.test, func(t *testing.T) {
				e := scenario.bag.Populate(scenario.target)

				if scenario.expectedErr != nil {
					assert.ErrorIs(t, e, scenario.expectedErr)
					return
				}

				assert.NoError(t, e)
				assert.Equal(t, scenario.expected, scenario.target)
			})
		}
	})

	t.Run("with path", func(t *testing.T) {
		scenarios := []struct {
			test        string
			bag         Bag
			path        string
			target      any
			expected    any
			expectedErr error
		}{
			{
				test:     "should populate a struct from a nested path",
				bag:      Bag{"config": Bag{"field": 456}},
				path:     "config",
				target:   &simpleStruct{},
				expected: &simpleStruct{Field: 456},
			},
			{
				test:        "should return an error for an invalid path",
				bag:         Bag{"config": Bag{"field": 456}},
				path:        "invalid.path",
				target:      &simpleStruct{},
				expectedErr: ErrBagInvalidPath,
			},
			{
				test: "should populate a complex struct from a nested path",
				bag: Bag{
					"data": Bag{
						"name":  "nested_name",
						"value": 111,
						"Nested": Bag{
							"Field": 222,
						},
						"Children": []any{"c1"},
					},
				},
				path:   "data",
				target: &complexStruct{},
				expected: &complexStruct{
					Name:     "nested_name",
					Value:    111,
					Nested:   simpleStruct{Field: 222},
					Children: []string{"c1"},
				},
			},
		}

		for _, scenario := range scenarios {
			t.Run(scenario.test, func(t *testing.T) {
				e := scenario.bag.Populate(scenario.target, scenario.path)

				if scenario.expectedErr != nil {
					assert.ErrorIs(t, e, scenario.expectedErr)
					return
				}

				assert.NoError(t, e)
				assert.Equal(t, scenario.expected, scenario.target)
			})
		}
	})
}
