package varint

import (
	"errors"
	"math/bits"
	"slices"
)

var (
	ErrNoBytes           = errors.New("varint: no bytes")
	ErrMissingTerminator = errors.New("varint: missing terminator")
	ErrTooLarge          = errors.New("varint: value too large for int")
)

func Decode(bytes []byte) (value int, length int, err error) {
	if len(bytes) == 0 {
		return 0, 0, ErrNoBytes
	}

	value = -1

	for _, b := range bytes {
		value += 1
		value <<= 7
		value |= int(b & 0b0111_1111)

		length += 1

		if b>>7 == 0b0 {
			if 7*length > bits.UintSize-1 {
				return 0, 0, ErrTooLarge
			}
			return value, length, nil
		}
	}

	return 0, 0, ErrMissingTerminator
}

func Encode(value int) []byte {
	bytes := []byte{}

	bytes = append(bytes, byte(value)&0b0111_1111)
	value >>= 7

	for value > 0 {
		value -= 1
		bytes = append(bytes, 0b1000_0000|byte(value))
		value >>= 7
	}

	slices.Reverse(bytes)

	return bytes
}
