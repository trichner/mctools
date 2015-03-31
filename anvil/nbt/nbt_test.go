// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package nbt

import (
	"bytes"
	"compress/gzip"
	"reflect"
	"testing"
	"time"
)

func TestByte1(t *testing.T) {
	var a, b uint8
	a = 123
	testRoundtrip(t, &a, &b)
}

func TestByte2(t *testing.T) {
	var a, b int8
	a = -123
	testRoundtrip(t, &a, &b)
}

func TestByte3(t *testing.T) {
	var a, b bool
	a = true
	testRoundtrip(t, &a, &b)
}

func TestShort1(t *testing.T) {
	var a, b uint16
	a = 123
	testRoundtrip(t, &a, &b)
}

func TestShort2(t *testing.T) {
	var a, b int16
	a = -123
	testRoundtrip(t, &a, &b)
}

func TestInt1(t *testing.T) {
	var a, b uint32
	a = 123
	testRoundtrip(t, &a, &b)
}

func TestInt2(t *testing.T) {
	var a, b int32
	a = -123
	testRoundtrip(t, &a, &b)
}

func TestLong1(t *testing.T) {
	var a, b uint64
	a = 123
	testRoundtrip(t, &a, &b)
}

func TestLong2(t *testing.T) {
	var a, b int64
	a = -123
	testRoundtrip(t, &a, &b)
}

func TestString1(t *testing.T) {
	var a, b string
	a = "test string"
	testRoundtrip(t, &a, &b)
}

func TestString2(t *testing.T) {
	var a, b string
	testRoundtrip(t, &a, &b)
}

func TestFloat1(t *testing.T) {
	var a, b float32
	a = 1.234
	testRoundtrip(t, &a, &b)
}

func TestFloat2(t *testing.T) {
	var a, b float32
	a = -1.234
	testRoundtrip(t, &a, &b)
}

func TestDouble1(t *testing.T) {
	var a, b float64
	a = 1.234
	testRoundtrip(t, &a, &b)
}

func TestDouble2(t *testing.T) {
	var a, b float64
	a = -1.234
	testRoundtrip(t, &a, &b)
}

func TestTime(t *testing.T) {
	type Test struct {
		Data time.Time
	}

	var a, b Test
	a.Data = time.Now().Truncate(time.Second)
	testRoundtrip(t, &a, &b)
}

func TestCompound(t *testing.T) {
	type Test struct {
		A struct {
			A int8
			B string
		}
	}

	var a, b Test
	a.A.A = 123
	a.A.B = "test"
	testRoundtrip(t, &a, &b)
}

func TestCompoundEmbedded(t *testing.T) {
	type T struct {
		A int8
		B string
	}

	type Test struct {
		T
	}

	var a, b Test
	a.A = 123
	a.B = "test"
	testRoundtrip(t, &a, &b)
}

func TestCompoundPtr1(t *testing.T) {
	type A struct {
		A int8
		B string
	}

	type Test struct {
		A *A
	}

	var a, b Test
	a.A = &A{
		A: 123,
		B: "test",
	}

	testRoundtrip(t, &a, &b)
}

func TestCompoundPtr2(t *testing.T) {
	type A struct {
		A int8
		B string
	}

	type Test struct {
		A *A
		C int32
	}

	var a, b Test
	a.A = &A{
		A: 123,
		B: "test",
	}
	a.C = -321

	testRoundtrip(t, &a, &b)
}

func TestCompoundPtr3(t *testing.T) {
	type A struct {
		A int8
		B string
	}

	type Test struct {
		A *A
		C int32
	}

	var a, b Test
	a.C = -321

	testRoundtrip(t, &a, &b)
}

func TestList1(t *testing.T) {
	var a, b []float32
	testRoundtrip(t, &a, &b)
}

func TestList2(t *testing.T) {
	type MyInt uint64

	var a, b []MyInt
	a = []MyInt{1, 2, 3, 4, 5}
	testRoundtrip(t, &a, &b)
}

func TestCompoundList(t *testing.T) {
	type Data struct {
		A int8
		B string
		C []int32
	}

	var a, b []Data
	a = []Data{
		{1, "test", []int32{1, 2, 3}},
		{2, "test", []int32{1, 2, 3}},
		{3, "test", []int32{1, 2, 3}},
		{4, "test", []int32{1, 2, 3}},
	}

	testRoundtrip(t, &a, &b)
}

func TestCompoundPtrList(t *testing.T) {
	type Data struct {
		A int8
		B string
		C []int32
	}

	var a, b []*Data
	a = []*Data{
		&Data{1, "test", []int32{1, 2, 3}},
		&Data{2, "test", []int32{1, 2, 3}},
		&Data{3, "test", []int32{1, 2, 3}},
		&Data{4, "test", []int32{1, 2, 3}},
	}

	testRoundtrip(t, &a, &b)
}

func TestBig(t *testing.T) {
	var a, b BigTest
	load(t, big_nbt, &a)
	testRoundtrip(t, &a, &b)
}

func TestSmall(t *testing.T) {
	var a, b SmallTest
	load(t, small_nbt, &a)
	testRoundtrip(t, &a, &b)
}

// testRoundtrip encodes <want> and then decodes into <have>.
// The two should then be equal.
func testRoundtrip(t *testing.T, want, have interface{}) {
	var buf bytes.Buffer

	enc := NewEncoder(&buf)
	err := enc.Encode(want)
	if err != nil {
		t.Fatal(err)
	}

	dec := NewDecoder(&buf)
	err = dec.Decode(have)
	if err != nil {
		t.Fatal(err)
	}

	rw := reflect.ValueOf(want)
	rw = reflect.Indirect(rw)
	want = rw.Interface()

	rh := reflect.ValueOf(have)
	rh = reflect.Indirect(rh)
	have = rh.Interface()

	if !reflect.DeepEqual(want, have) {
		t.Fatalf("roundtrip mismatch:\nHave: %#v\nWant: %#v", have, want)
	}
}

type Food struct {
	Name  string  `nbt:"name"`
	Value float32 `nbt:"value"`
}

type NestedCompound struct {
	Egg Food `nbt:"egg"`
	Ham Food `nbt:"ham"`
}

type ListItem struct {
	Name    string `nbt:"name"`
	Created int64  `nbt:"created-on"`
}

type BigTestData struct {
	NestedCompound `nbt:"nested compound test"`
	ListTestCmp    []ListItem `nbt:"listTest (compound)"`
	ListTestLong   []int64    `nbt:"listTest (long)"`
	ListTestByte   []byte     `nbt:"byteArrayTest (the first 1000 values of (n*n*255+n*7)%100, starting with n=0 (0, 62, 34, 16, 8, ...))"`
	StringTest     string     `nbt:"stringTest"`
	LongTest       int64      `nbt:"longTest"`
	DoubleTest     float64    `nbt:"doubleTest"`
	FloatTest      float32    `nbt:"floatTest"`
	IntTest        int32      `nbt:"intTest"`
	ShortTest      int16      `nbt:"shortTest"`
	ByteTest       byte       `nbt:"byteTest"`
}

type BigTest struct {
	BigTestData `nbt:"Level"`
}

type SmallTestData struct {
	Name string `nbt:"name"`
}

type SmallTest struct {
	SmallTestData `nbt:"hello world"`
}

func load(t *testing.T, data []byte, v interface{}) {
	r, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}

	defer r.Close()

	err = Unmarshal(r, v)
	if err != nil {
		t.Fatal(err)
	}
}

var big_nbt = []byte{
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xed, 0x54,
	0xcf, 0x4f, 0x1a, 0x41, 0x14, 0x7e, 0xc2, 0x02, 0xcb, 0x96, 0x82, 0xb1,
	0xc4, 0x10, 0x63, 0xcc, 0xab, 0xb5, 0x84, 0xa5, 0xdb, 0xcd, 0x42, 0x11,
	0x89, 0xb1, 0x88, 0x16, 0x2c, 0x9a, 0x0d, 0x1a, 0xd8, 0xa8, 0x31, 0x86,
	0xb8, 0x2b, 0xc3, 0x82, 0x2e, 0xbb, 0x66, 0x77, 0xb0, 0xf1, 0xd4, 0x4b,
	0x7b, 0x6c, 0x7a, 0xeb, 0x3f, 0xd3, 0x23, 0x7f, 0x43, 0xcf, 0xbd, 0xf6,
	0xbf, 0xa0, 0xc3, 0x2f, 0x7b, 0x69, 0xcf, 0xbd, 0xf0, 0x32, 0xc9, 0xf7,
	0xe6, 0xbd, 0x6f, 0xe6, 0x7b, 0x6f, 0x26, 0x79, 0x02, 0x04, 0x54, 0x72,
	0x4f, 0x2c, 0x0e, 0x78, 0xcb, 0xb1, 0x4d, 0x8d, 0x78, 0xf4, 0xe3, 0x70,
	0x62, 0x3e, 0x08, 0x7b, 0x1d, 0xc7, 0xa5, 0x93, 0x18, 0x0f, 0x82, 0x47,
	0xdd, 0xee, 0x84, 0x02, 0x62, 0xb5, 0xa2, 0xaa, 0xc7, 0x78, 0x76, 0x5c,
	0x57, 0xcb, 0xa8, 0x55, 0x0f, 0x1b, 0xc8, 0xd6, 0x1e, 0x6a, 0x95, 0x86,
	0x86, 0x0d, 0xad, 0x7e, 0x58, 0x7b, 0x8f, 0x83, 0xcf, 0x83, 0x4f, 0x83,
	0x6f, 0xcf, 0x03, 0x10, 0x6e, 0x5b, 0x8e, 0x3e, 0xbe, 0xa5, 0x38, 0x4c,
	0x64, 0xfd, 0x10, 0xea, 0xda, 0x74, 0xa6, 0x23, 0x40, 0xdc, 0x66, 0x2e,
	0x69, 0xe1, 0xb5, 0xd3, 0xbb, 0x73, 0xfa, 0x76, 0x0b, 0x29, 0xdb, 0x0b,
	0xe0, 0xef, 0xe8, 0x3d, 0x1e, 0x38, 0x5b, 0xef, 0x11, 0x08, 0x56, 0xf5,
	0xde, 0x5d, 0xdf, 0x0b, 0x40, 0xe0, 0x5e, 0xb7, 0xfa, 0x64, 0xb7, 0x04,
	0x00, 0x8c, 0x41, 0x4c, 0x73, 0xc6, 0x08, 0x55, 0x4c, 0xd3, 0x20, 0x2e,
	0x7d, 0xa4, 0xc0, 0xc8, 0xc2, 0x10, 0xb3, 0xba, 0xde, 0x58, 0x0b, 0x53,
	0xa3, 0xee, 0x44, 0x8e, 0x45, 0x03, 0x30, 0xb1, 0x27, 0x53, 0x8c, 0x4c,
	0xf1, 0xe9, 0x14, 0xa3, 0x53, 0x8c, 0x85, 0xe1, 0xd9, 0x9f, 0xe3, 0xb3,
	0xf2, 0x44, 0x81, 0xa5, 0x7c, 0x33, 0xdd, 0xd8, 0xbb, 0xc7, 0xaa, 0x75,
	0x13, 0x5f, 0x28, 0x1c, 0x08, 0xd7, 0x2e, 0xd1, 0x59, 0x3f, 0xaf, 0x1d,
	0x1b, 0x60, 0x21, 0x59, 0xdf, 0xfa, 0xf1, 0x05, 0xfe, 0xc1, 0xce, 0xfc,
	0x9d, 0xbd, 0x00, 0xbc, 0xf1, 0x40, 0xc9, 0xf8, 0x85, 0x42, 0x40, 0x46,
	0xfe, 0x9e, 0xeb, 0xea, 0x0f, 0x93, 0x3a, 0x68, 0x87, 0x60, 0xbb, 0xeb,
	0x32, 0x37, 0xa3, 0x28, 0x0a, 0x8e, 0xbb, 0xf5, 0xd0, 0x69, 0x63, 0xca,
	0x4e, 0xdb, 0xe9, 0xec, 0xe6, 0xe6, 0x2b, 0x3b, 0xbd, 0x25, 0xbe, 0x64,
	0x49, 0x09, 0x3d, 0xaa, 0xbb, 0x94, 0xfd, 0x18, 0x7e, 0xe8, 0xd2, 0x0e,
	0xda, 0x6f, 0x15, 0x4c, 0xb1, 0x68, 0x3e, 0x2b, 0xe1, 0x9b, 0x9c, 0x84,
	0x99, 0xbc, 0x84, 0x05, 0x09, 0x65, 0x59, 0x16, 0x45, 0x00, 0xff, 0x2f,
	0x28, 0xae, 0x2f, 0xf2, 0xc2, 0xb2, 0xa4, 0x2e, 0x1d, 0x20, 0x77, 0x5a,
	0x3b, 0xb9, 0x8c, 0xca, 0xe7, 0x29, 0xdf, 0x51, 0x41, 0xc9, 0x16, 0xb5,
	0xc5, 0x6d, 0xa1, 0x2a, 0xad, 0x2c, 0xc5, 0x31, 0x7f, 0xba, 0x7a, 0x92,
	0x8e, 0x5e, 0x9d, 0x5f, 0xf8, 0x12, 0x05, 0x23, 0x1b, 0xd1, 0xf6, 0xb7,
	0x77, 0xaa, 0xcd, 0x95, 0x72, 0xbc, 0x9e, 0xdf, 0x58, 0x5d, 0x4b, 0x97,
	0xae, 0x92, 0x17, 0xb9, 0x44, 0xd0, 0x80, 0xc8, 0xfa, 0x3e, 0xbf, 0xb3,
	0xdc, 0x54, 0xcb, 0x07, 0x75, 0x6e, 0xa3, 0xb6, 0x76, 0x59, 0x92, 0x93,
	0xa9, 0xdc, 0x51, 0x50, 0x99, 0x6b, 0xcc, 0x35, 0xe6, 0x1a, 0xff, 0x57,
	0x23, 0x08, 0x42, 0xcb, 0xe9, 0x1b, 0xd6, 0x78, 0xc2, 0xec, 0xfe, 0xfc,
	0x7a, 0xfb, 0x7d, 0x78, 0xd3, 0x84, 0xdf, 0xd4, 0xf2, 0xa4, 0xfb, 0x08,
	0x06, 0x00, 0x00,
}

var small_nbt = []byte{
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xe3, 0x62,
	0xe0, 0xce, 0x48, 0xcd, 0xc9, 0xc9, 0x57, 0x28, 0xcf, 0x2f, 0xca, 0x49,
	0xe1, 0x60, 0x60, 0xc9, 0x4b, 0xcc, 0x4d, 0x65, 0xe0, 0x74, 0x4a, 0xcc,
	0x4b, 0xcc, 0x2b, 0x4a, 0xcc, 0x4d, 0x64, 0x00, 0x00, 0x77, 0xda, 0x5c,
	0x3a, 0x21, 0x00, 0x00, 0x00,
}