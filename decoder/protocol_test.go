package decoder

import (
	"bytes"
	"testing"
)

// StartDatabase
var selecttests = []struct {
	db       int
	expected string
}{
	{1, "*2\r\n$6\r\nSELECT\r\n$1\r\n1\r\n"},
	{10, "*2\r\n$6\r\nSELECT\r\n$2\r\n10\r\n"},
	{435, "*2\r\n$6\r\nSELECT\r\n$3\r\n435\r\n"},
}

func TestStartDatabase(t *testing.T) {
	var buf bytes.Buffer
	d := NewProtocolDecoder(&buf)

	for _, tt := range selecttests {
		buf.Reset()
		d.StartDatabase(tt.db)
		if buf.String() != tt.expected {
			t.Errorf("StartDatabase(%d) => %q, want %q", tt.db, buf.String(), tt.expected)
		}
	}

}

// Set
var settests = []struct {
	key      []byte
	value    []byte
	expire   int64
	expected string
}{
	{[]byte("ala"), []byte("ma kota"), 1671963072573, "*3\r\n$3\r\nSET\r\n$3\r\nala\r\n$7\r\nma kota\r\n*3\r\n$9\r\nPEXPIREAT\r\n$3\r\nala\r\n$13\r\n1671963072573\r\n"},
	{[]byte("pokemon"), []byte("pawel"), 0, "*3\r\n$3\r\nSET\r\n$7\r\npokemon\r\n$5\r\npawel\r\n"},
}

func TestSet(t *testing.T) {
	var buf bytes.Buffer
	d := NewProtocolDecoder(&buf)

	for _, tt := range settests {
		buf.Reset()
		d.Set(tt.key, tt.value, tt.expire)
		if buf.String() != tt.expected {
			t.Errorf("Set(%q, %q, %d) => %q, want %q", tt.key, tt.value, tt.expire, buf.String(), tt.expected)
		}
	}
}

// HSet
var hsettests = []struct {
	key      []byte
	field    []byte
	value    []byte
	expire   int64
	expected string
}{
	{[]byte("foo"), []byte("bar"), []byte("baz"), 1671963072573, "*4\r\n$4\r\nHSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n$3\r\nbaz\r\n*3\r\n$9\r\nPEXPIREAT\r\n$3\r\nfoo\r\n$13\r\n1671963072573\r\n"},
	{[]byte("foo"), []byte("bar"), []byte("baz"), 12345, "*4\r\n$4\r\nHSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n$3\r\nbaz\r\n*3\r\n$9\r\nPEXPIREAT\r\n$3\r\nfoo\r\n$5\r\n12345\r\n"},
	{[]byte("catch"), []byte("em"), []byte("all"), 0, "*4\r\n$4\r\nHSET\r\n$5\r\ncatch\r\n$2\r\nem\r\n$3\r\nall\r\n"},
}

func TestHSet(t *testing.T) {
	var buf bytes.Buffer
	d := NewProtocolDecoder(&buf)

	for _, tt := range hsettests {
		buf.Reset()

		d.StartHash(tt.key, 4, tt.expire)
		d.Hset(tt.key, tt.field, tt.value)
		d.EndHash(tt.key)

		if buf.String() != tt.expected {
			t.Errorf("Hset(%q, %q, %q, %d) => %q, want %q", tt.key, tt.field, tt.value, tt.expire, buf.String(), tt.expected)
		}
	}
}
