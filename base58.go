package base58

import (
	"errors"
	"unicode"
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

var ErrInvalidAlphabet = errors.New("invalid base58 alphabet")

const (
	Flickr  = "123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"
	Ripple  = "rpshnaf39wBUDNEGHJKLM4PQRST7VWXYZ2bcdeCg65jkm8oFqi1tuvAxyz"
	Bitcoin = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
)

func NewEncodeDecoder(alphabet string, offset int64) (EncodeDecoder, error) {
	if !isASCII(alphabet) {
		return nil, ErrInvalidAlphabet
	}

	if len(alphabet) != 58 {
		return nil, ErrInvalidAlphabet
	}

	if offset < 0 {
		offset = 0
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

func (ed encodeDecoder) Encode(id int64) string {
	if id == 0 {
		return ""
	}
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

func (ed encodeDecoder) Decode(s string) int64 {
	if s == "" {
		return 0
	}
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
