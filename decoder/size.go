package decoder

import (
	"fmt"
	"github.com/cupcake/rdb/nopdecoder"
)

type sizeAttr struct {
	key   int
	hash  int
	set   int
	zset  int
	list  int
	bytes int
}

type size struct {
	db int
	sizeAttr
	exp sizeAttr
	nopdecoder.NopDecoder
}

func Size() *size {
	return &size{}
}

func (p *size) StartDatabase(n int) {
	if n > 0 {
		p.printEnd()
	}
	p.sizeAttr = sizeAttr{}
	p.exp = sizeAttr{}
	p.db = n
}

func pct(num, denom int) int {
	if denom == 0 {
		return 0
	}
	return num * 100 / denom
}

func (p *size) printEnd() {
	fmt.Printf("db=%2d  bytes=%10d  key=%10d %3d%%  hash=%10d %3d%%  set=%10d %3d%%  zset=%10d %3d%%  list=%10d %3d%%\n",
		p.db, p.bytes,
		p.key, pct(p.exp.key, p.key),
		p.hash, pct(p.exp.hash, p.hash),
		p.set, pct(p.exp.set, p.set),
		p.zset, pct(p.exp.zset, p.zset),
		p.list, pct(p.exp.list, p.list),
	)
}

func (p *size) EndDatabase(n int) { p.printEnd() }

func (p *size) Set(key, value []byte, expiry int64) {
	p.key++
	p.bytes += len(key) + len(value)
	if expiry > 0 {
		p.exp.key++
	}
}

func (p *size) StartHash(key []byte, length, expiry int64) {
	p.bytes += len(key)
	if expiry > 0 {
		p.exp.hash += int(length)
	}
}

func (p *size) Hset(key, field, value []byte) {
	p.hash++
	p.bytes += len(field) + len(value)
}

func (p *size) StartSet(key []byte, cardinality, expiry int64) {
	p.bytes += len(key)
	p.set += int(cardinality)
	if expiry > 0 {
		p.exp.set += int(cardinality)
	}
}

func (p *size) Sadd(key, member []byte) {
	p.bytes += len(member)
}

func (p *size) StartList(key []byte, length, expiry int64) {
	p.bytes += len(key)
	p.list += int(length)
	if expiry > 0 {
		p.exp.list += int(length)
	}
}

func (p *size) Rpush(key, value []byte) {
	p.bytes += len(value)
}

func (p *size) StartZSet(key []byte, cardinality, expiry int64) {
	p.bytes += len(key)
	p.zset += int(cardinality)
	if expiry > 0 {
		p.exp.zset += int(cardinality)
	}
}

func (p *size) Zadd(key []byte, score float64, member []byte) {
	p.bytes += len(member)
}
