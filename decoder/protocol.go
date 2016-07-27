package decoder

import (
	"bufio"
	"io"
	"strconv"
)

const NoExpire = 0

var (
	arrayPrefixSlice      = []byte{'*'}
	bulkStringPrefixSlice = []byte{'$'}
	lineEndingSlice       = []byte{'\r', '\n'}
)

type protocol struct {
	expire int64
	*bufio.Writer
}

func NewProtocolDecoder(writer io.Writer) *protocol {
	return &protocol{
		Writer: bufio.NewWriter(writer),
	}
}

func (p *protocol) writeCommand(args ...string) (err error) {
	// Write the array prefix and the number of arguments in the array.
	p.Write(arrayPrefixSlice)
	p.WriteString(strconv.Itoa(len(args)))
	p.Write(lineEndingSlice)

	// Write a bulk string for each argument.
	for _, arg := range args {
		p.Write(bulkStringPrefixSlice)
		p.WriteString(strconv.Itoa(len(arg)))
		p.Write(lineEndingSlice)
		p.WriteString(arg)
		p.Write(lineEndingSlice)
	}

	return p.Flush()
}

func (p *protocol) preExpire(expire int64) {
	p.expire = expire
}

func (p *protocol) postExpire(key string) {
	if p.expire == NoExpire {
		return
	}
	p.writeCommand("PEXPIREAT", key, strconv.FormatInt(p.expire, 10))
	p.expire = NoExpire
}

func (p *protocol) StartRDB() {}
func (p *protocol) EndRDB()   {}

func (p *protocol) StartDatabase(n int) {
	p.writeCommand("SELECT", strconv.Itoa(n))
}
func (p *protocol) EndDatabase(n int) {}

func (p *protocol) Set(key, value []byte, expire int64) {
	p.preExpire(expire)
	p.writeCommand("SET", string(key), string(value))
	p.postExpire(string(key))
}

func (p *protocol) StartHash(key []byte, length, expire int64) {
	p.preExpire(expire)
}

func (p *protocol) Hset(key, field, value []byte) {
	p.writeCommand("HSET", string(key), string(field), string(value))
}

func (p *protocol) EndHash(key []byte) {
	p.postExpire(string(key))
}

func (p *protocol) StartSet(key []byte, cardinality, expire int64) {
	p.preExpire(expire)
}

func (p *protocol) Sadd(key, member []byte) {
	p.writeCommand("SADD", string(key), string(member))
}

func (p *protocol) EndSet(key []byte) {
	p.postExpire(string(key))
}

func (p *protocol) StartList(key []byte, length, expire int64) {
	p.preExpire(expire)
}
func (p *protocol) Rpush(key, value []byte) {
	p.writeCommand("RPUSH", string(key), string(value))
}
func (p *protocol) EndList(key []byte) {
	p.postExpire(string(key))
}

func (p *protocol) StartZSet(key []byte, cardinality, expire int64) {
	p.preExpire(expire)
}
func (p *protocol) Zadd(key []byte, score float64, member []byte) {
	p.writeCommand("ZADD", string(key), strconv.FormatFloat(score, 'f', 6, 64), string(member))
}
func (p *protocol) EndZSet(key []byte) {
	p.postExpire(string(key))
}
