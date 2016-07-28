package decoder

import (
	"bytes"
	"testing"
)

// Database selection
var selecttests = []struct {
	db       int
	expected string
}{
	{1, "*2\r\n$6\r\nSELECT\r\n$1\r\n1\r\n"},
	{10, "*2\r\n$6\r\nSELECT\r\n$2\r\n10\r\n"},
	{435, "*2\r\n$6\r\nSELECT\r\n$3\r\n435\r\n"},
}

func TestDatabaseSelect(t *testing.T) {
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

// Simple KV using Set
var settests = []struct {
	key      []byte
	value    []byte
	expire   int64
	expected string
}{
	{[]byte("ala"), []byte("ma kota"), 1671963072573, "*3\r\n$3\r\nSET\r\n$3\r\nala\r\n$7\r\nma kota\r\n*3\r\n$9\r\nPEXPIREAT\r\n$3\r\nala\r\n$13\r\n1671963072573\r\n"},
	{[]byte("pokemon"), []byte("pawel"), 0, "*3\r\n$3\r\nSET\r\n$7\r\npokemon\r\n$5\r\npawel\r\n"},
}

func TestSimpleKvSet(t *testing.T) {
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

// Hash maps
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

func TestHashMaps(t *testing.T) {
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

// Sets
var saddtests = []struct {
	key      []byte
	member   []byte
	expire   int64
	expected string
}{
	{[]byte("ala"), []byte("kot"), 908761, "*3\r\n$4\r\nSADD\r\n$3\r\nala\r\n$3\r\nkot\r\n*3\r\n$9\r\nPEXPIREAT\r\n$3\r\nala\r\n$6\r\n908761\r\n"},
	{[]byte("ala"), []byte("pies"), 3451, "*3\r\n$4\r\nSADD\r\n$3\r\nala\r\n$4\r\npies\r\n*3\r\n$9\r\nPEXPIREAT\r\n$3\r\nala\r\n$4\r\n3451\r\n"},
	{[]byte("pokemon"), []byte("pawel"), 0, "*3\r\n$4\r\nSADD\r\n$7\r\npokemon\r\n$5\r\npawel\r\n"},
	{[]byte("pokemon"), []byte("pikachu"), 0, "*3\r\n$4\r\nSADD\r\n$7\r\npokemon\r\n$7\r\npikachu\r\n"},
}

func TestSets(t *testing.T) {
	var buf bytes.Buffer
	d := NewProtocolDecoder(&buf)

	for _, tt := range saddtests {
		buf.Reset()

		d.StartSet(tt.key, 4, tt.expire)
		d.Sadd(tt.key, tt.member)
		d.EndHash(tt.key)

		if buf.String() != tt.expected {
			t.Errorf("Sadd(%q, %q) => %q, want %q", tt.key, tt.member, buf.String(), tt.expected)
		}
	}
}

// Lists
var listtests = []struct {
	key      []byte
	value    []byte
	expire   int64
	expected string
}{
	{[]byte("ala"), []byte("kot"), 908761, "*3\r\n$5\r\nRPUSH\r\n$3\r\nala\r\n$3\r\nkot\r\n*3\r\n$9\r\nPEXPIREAT\r\n$3\r\nala\r\n$6\r\n908761\r\n"},
	{[]byte("ala"), []byte("pies"), 3451, "*3\r\n$5\r\nRPUSH\r\n$3\r\nala\r\n$4\r\npies\r\n*3\r\n$9\r\nPEXPIREAT\r\n$3\r\nala\r\n$4\r\n3451\r\n"},
	{[]byte("pokemon"), []byte("pawel"), 0, "*3\r\n$5\r\nRPUSH\r\n$7\r\npokemon\r\n$5\r\npawel\r\n"},
	{[]byte("pokemon"), []byte("pikachu"), 0, "*3\r\n$5\r\nRPUSH\r\n$7\r\npokemon\r\n$7\r\npikachu\r\n"},
}

func TestLists(t *testing.T) {
	var buf bytes.Buffer
	d := NewProtocolDecoder(&buf)

	for _, tt := range listtests {
		buf.Reset()

		d.StartList(tt.key, 1, tt.expire)
		d.Rpush(tt.key, tt.value)
		d.EndList(tt.key)

		if buf.String() != tt.expected {
			t.Errorf("Rpush(%q, %q) => %q, want %q", tt.key, tt.value, buf.String(), tt.expected)
		}
	}
}

// Sorted Sets
var sortedsettests = []struct {
	key      []byte
	score    float64
	member   []byte
	expire   int64
	expected string
}{
	{[]byte("ala"), 1.2, []byte("kot"), 123456, "*4\r\n$4\r\nZADD\r\n$3\r\nala\r\n$3\r\n1.2\r\n$3\r\nkot\r\n*3\r\n$9\r\nPEXPIREAT\r\n$3\r\nala\r\n$6\r\n123456\r\n"},
	{[]byte("foo"), 2, []byte("barbaz"), 123456, "*4\r\n$4\r\nZADD\r\n$3\r\nfoo\r\n$1\r\n2\r\n$6\r\nbarbaz\r\n*3\r\n$9\r\nPEXPIREAT\r\n$3\r\nfoo\r\n$6\r\n123456\r\n"},
	{[]byte("ala"), 2.34, []byte("kot"), 123456, "*4\r\n$4\r\nZADD\r\n$3\r\nala\r\n$4\r\n2.34\r\n$3\r\nkot\r\n*3\r\n$9\r\nPEXPIREAT\r\n$3\r\nala\r\n$6\r\n123456\r\n"},
}

func TestSortedSets(t *testing.T) {
	var buf bytes.Buffer
	d := NewProtocolDecoder(&buf)

	for _, tt := range sortedsettests {
		buf.Reset()

		d.StartZSet(tt.key, 1, tt.expire)
		d.Zadd(tt.key, tt.score, tt.member)
		d.EndZSet(tt.key)

		if buf.String() != tt.expected {
			t.Errorf("Zadd(%q, %f, %q) => %q, want %q", tt.key, tt.score, tt.member, buf.String(), tt.expected)
		}
	}
}
