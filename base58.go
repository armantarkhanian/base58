package base58

import (
	"errors"
	"unicode"
	"unicode/utf8"
)

type EncodeDecoder interface {
	Encode(int64) string
	Decode(string) int64
}

type encodeDecoder struct {
	offset          int64
	alphabet        string
	decodeBase58Map [256]byte
}

var _ EncodeDecoder = &encodeDecoder{}

var (
	ErrInvalidAlphabet = errors.New("base58: invalid alphabet")
	ErrInvalidOffset   = errors.New("base58: invalid offset [0; 1 000 000 000]")
)

const (
	MaxOffset = 1000000000
)

const (
	Flickr  = "123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"
	Ripple  = "rpshnaf39wBUDNEGHJKLM4PQRST7VWXYZ2bcdeCg65jkm8oFqi1tuvAxyz"
	Bitcoin = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
)

func NewEncodeDecoder(alphabet string, offset int64) (EncodeDecoder, error) {
	if !isASCII(alphabet) {
		return nil, ErrInvalidAlphabet
	}

	if utf8.RuneCountInString(alphabet) != 58 {
		return nil, ErrInvalidAlphabet
	}

	if offset < 0 || offset > MaxOffset {
		return nil, ErrInvalidOffset
	}

	ed := encodeDecoder{
		alphabet: alphabet,
		offset:   offset,
	}

	for i := 0; i < len(ed.decodeBase58Map); i++ {
		ed.decodeBase58Map[i] = 0xFF
	}

	for i := 0; i < len(ed.alphabet); i++ {
		ed.decodeBase58Map[ed.alphabet[i]] = byte(i)
	}

	return &ed, nil
}

// No need to check error because we got id from database, which means it always should be correct
// and when we create EncodeDecoder interface we also check offset range, so MaxOffset is only 1000000000, so
// we have 9 223 372 036 854 775 807 - 1 000 000 000 = 9 223 372 035 854 775 807 possible integers
func (ed encodeDecoder) Encode(id int64) string {
	id += ed.offset
	if id < 58 {
		return string(ed.alphabet[id])
	}

	b := make([]byte, 0, 11)
	for id >= 58 {
		b = append(b, ed.alphabet[id%58])
		id /= 58
	}
	b = append(b, ed.alphabet[id])

	for x, y := 0, len(b)-1; x < y; x, y = x+1, y-1 {
		b[x], b[y] = b[y], b[x]
	}

	return string(b)
}

// No need to check error because decode will always return integer and for application it will just mean that client sent invalid ID
func (ed encodeDecoder) Decode(s string) int64 {
	b := []byte(s)

	var id int64

	for i := range b {
		id = id*58 + int64(ed.decodeBase58Map[b[i]])
	}

	id -= ed.offset

	return id
}

func isASCII(s string) bool {
	for _, c := range s {
		if c > unicode.MaxASCII {
			return false
		}
	}
	return true
}
