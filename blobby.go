package blobby

import (
	"errors"

	"github.com/nixberg/blobby-go/internal/varint"
)

var (
	ErrVarint                = errors.New("blobby: bad varint")
	ErrBlobLength            = errors.New("blobby: bad blob length")
	ErrDeduplicatedBlobIndex = errors.New("blobby: bad index for deduplicated blob")
)

func Decode(bytes []byte) ([][]byte, error) {
	popVarint := func() (int, error) {
		varint, length, err := varint.Decode(bytes)
		if err != nil {
			return 0, err
		}
		bytes = bytes[length:]
		return varint, nil
	}

	deduplicatedBlobsCount, err := popVarint()
	if err != nil {
		return nil, errors.Join(ErrVarint, err)
	}

	deduplicatedBlobs := make([][]byte, 0, deduplicatedBlobsCount)

	for i := 0; i < deduplicatedBlobsCount; i++ {
		length, err := popVarint()
		if err != nil {
			return nil, errors.Join(ErrVarint, err)
		}
		if length > len(bytes) {
			return nil, ErrBlobLength
		}
		deduplicatedBlobs = append(deduplicatedBlobs, bytes[:length])
		bytes = bytes[length:]
	}

	blobs := [][]byte{}

	for len(bytes) > 0 {
		indexOrLength, err := popVarint()
		if err != nil {
			return nil, errors.Join(ErrVarint, err)
		}

		if indexOrLength&0b1 == 1 {
			index := indexOrLength >> 1
			if index >= deduplicatedBlobsCount {
				return nil, ErrDeduplicatedBlobIndex
			}
			blobs = append(blobs, deduplicatedBlobs[index])
		} else {
			length := indexOrLength >> 1
			if length > len(bytes) {
				return nil, ErrBlobLength
			}
			blobs = append(blobs, bytes[:length])
			bytes = bytes[length:]
		}
	}

	return blobs, nil
}

func MustDecode(bytes []byte) [][]byte {
	result, err := Decode(bytes)
	if err != nil {
		panic(err)
	}
	return result
}
