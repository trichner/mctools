// generated by stringer -type=TagId; DO NOT EDIT

package main

import "fmt"

const (
	_TagId_name_0 = "TagEndTagByteTagShortTagIntTagLongTagFloatTagDoubleTagByteArrayTagStringTagListTagCompoundTagIntArray"
	_TagId_name_1 = "TagUnknown"
)

var (
	_TagId_index_0 = [...]uint8{6, 13, 21, 27, 34, 42, 51, 63, 72, 79, 90, 101}
	_TagId_index_1 = [...]uint8{10}
)

func (i tagId) String() string {
	switch {
	case 0 <= i && i <= 11:
		lo := uint8(0)
		if i > 0 {
			lo = _TagId_index_0[i-1]
		}
		return _TagId_name_0[lo:_TagId_index_0[i]]
	case i == 255:
		return _TagId_name_1
	default:
		return fmt.Sprintf("TagId(%d)", i)
	}
}
