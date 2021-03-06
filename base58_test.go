package base58

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncode(t *testing.T) {
	tt := []struct {
		integer        int64
		offset         int64
		expectedBase58 string
		valid          bool
	}{
		{
			integer:        15,
			offset:         0,
			expectedBase58: "g",
			valid:          true,
		},
		{
			integer:        0,
			offset:         0,
			expectedBase58: "1",
			valid:          true,
		},
		{
			integer:        15,
			offset:         1,
			expectedBase58: "h",
			valid:          true,
		},
		{
			integer:        15,
			expectedBase58: "w",
			valid:          false,
		},
		{
			integer:        59,
			expectedBase58: "22",
			valid:          true,
		},
	}

	for _, tc := range tt {
		ed, err := New(Flickr, tc.offset)
		require.NoError(t, err)
		require.NotNil(t, ed)
		base58String := ed.Encode(tc.integer)
		if tc.valid {
			require.Equal(t, tc.expectedBase58, base58String)
		} else {
			require.NotEqual(t, tc.expectedBase58, base58String)
		}
	}
}

func TestDecode(t *testing.T) {
	tt := []struct {
		base58          string
		offset          int64
		expectedInteger int64
		valid           bool
	}{
		{
			base58:          "",
			offset:          0,
			expectedInteger: 0,
			valid:           true,
		},
		{
			base58:          "g",
			offset:          0,
			expectedInteger: 15,
			valid:           true,
		},
		{
			base58:          "h",
			offset:          1,
			expectedInteger: 15,
			valid:           true,
		},
		{
			base58:          "w",
			offset:          0,
			expectedInteger: 15,
			valid:           false,
		},
	}

	for _, tc := range tt {
		ed, err := New(Flickr, tc.offset)
		require.NoError(t, err)
		require.NotNil(t, ed)
		integer := ed.Decode(tc.base58)
		if tc.valid {
			require.Equal(t, tc.expectedInteger, integer)
		} else {
			require.NotEqual(t, tc.expectedInteger, integer)
		}
	}
}

func TestNewEncodeDecoder(t *testing.T) {
	tt := []struct {
		alphabet string
		offset   int64
		valid    bool
	}{
		{
			alphabet: "",
			offset:   1,
			valid:    false,
		},
		{
			alphabet: "????????????",
			offset:   1,
			valid:    false,
		},
		{
			alphabet: Flickr,
			offset:   2,
			valid:    true,
		},
		{
			alphabet: Bitcoin,
			offset:   0,
			valid:    true,
		},
		{
			alphabet: Ripple,
			offset:   -1,
			valid:    false,
		},
		{
			alphabet: "abcdefg",
			offset:   0,
			valid:    false,
		},
	}
	for _, tc := range tt {
		ed, err := New(tc.alphabet, tc.offset)
		if !tc.valid {
			require.Error(t, err)
			require.Nil(t, ed)
			continue
		}
		require.NoError(t, err)
		require.NotNil(t, ed)

		if tc.offset <= 0 {
			require.Equal(t, int64(0), ed.offset)
		} else {
			require.Equal(t, tc.offset, ed.offset)
		}

		require.Equal(t, tc.alphabet, ed.alphabet)
		require.NotNil(t, ed.decodeBase58Map)
		require.Len(t, ed.decodeBase58Map, 256)
	}
}

func TestIsASCII(t *testing.T) {
	tt := []struct {
		s       string
		isASCII bool
	}{
		{"arman", true},
		{"admin", true},
		{"adm??in", false},
	}
	for _, tc := range tt {
		isASCII := isASCII(tc.s)
		require.Equal(t, tc.isASCII, isASCII)
	}
}
