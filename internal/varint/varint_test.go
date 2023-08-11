package varint

import (
	"bytes"
	"math/bits"
	"math/rand"
	"testing"
)

func TestDecodeValid(t *testing.T) {
	for _, vector := range []struct {
		bytes    []byte
		expected int
	}{
		{[]byte{0b0000_0000}, 0},
		{[]byte{0b0000_0010}, 2},
		{[]byte{0b0111_1111}, 127},
		{[]byte{0b1000_0000, 0b0000_0000}, 128},
		{[]byte{0b1111_1111, 0b0111_1111}, 16511},
		{[]byte{0b1000_0000, 0b1000_0000, 0b0000_0000}, 16512},
		{[]byte{0b1111_1111, 0b1111_1111, 0b0111_1111}, 2113663},
		{[]byte{0b1000_0000, 0b1000_0000, 0b1000_0000, 0b0000_0000}, 2113664},
		{[]byte{0b1111_1111, 0b1111_1111, 0b1111_1111, 0b0111_1111}, 270549119},
		{[]byte{0b1111_1111, 0b1111_1111, 0b1111_1111, 0b1111_1111, 0b0111_1111}, 34630287487},
	} {
		value, length, err := Decode(vector.bytes)
		if err != nil {
			t.Error(err)
		} else if length != len(vector.bytes) {
			t.Fail()
		} else if value != vector.expected {
			t.Fail()
		}

		if !bytes.Equal(Encode(value), vector.bytes) {
			t.Fail()
		}
	}
}

func TestDecodeMissingTerminator(t *testing.T) {
	_, _, err := Decode([]byte{})
	if err == nil {
		t.Fail()
	} else if err != ErrNoBytes {
		t.Errorf("unexpected error: %v", err)
	}
	for _, bytes := range [][]byte{
		{0b1000_0010},
		{0b1111_1111},
		{0b1111_1111, 0b1000_0000},
		{0b1111_1111, 0b1111_1111},
	} {
		_, _, err := Decode(bytes)
		if err == nil {
			t.Fail()
		} else if err != ErrMissingTerminator {
			t.Errorf("unexpected error: %v", err)
		}
	}
}

func TestDecodeTooLarge(t *testing.T) {
	if bits.UintSize > 64 {
		t.Fail()
	}
	_, _, err := Decode([]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0})
	if err == nil {
		t.Fail()
	} else if err != ErrTooLarge {
		t.Fail()
	}
}

func FuzzRoundtrip(f *testing.F) {
	for i := 0; i < 1024; i++ {
		f.Add(rand.Int())
	}
	f.Fuzz(func(t *testing.T, value int) {
		if value < 0 {
			t.Skip()
		}
		encoded := Encode(value)
		decoded, length, err := Decode(encoded)
		if err != nil {
			t.Fail()
		} else if length != len(encoded) {
			t.Fail()
		} else if decoded != value {
			t.Fail()
		}
	})
}
