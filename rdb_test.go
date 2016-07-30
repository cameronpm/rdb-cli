package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"github.com/brainly/rdb-cli/decoder"
	"testing"
)

var fixtures = []struct {
	name        string
	md5Protocol string
}{
	{"easily_compressible_string_key", "75425c72b9b7b8ddef70a00457732bc4"},
	{"empty_database", "d41d8cd98f00b204e9800998ecf8427e"},
	{"hash_as_ziplist", "7dccfe5b96da67f13aa6d42e29be04c6"},
	{"integer_keys", "e77bdc8792ad30e522db651730fb72f4"},
	{"intset_16", "b6f08218b2a80bf78ff4e0a3ba356b09"},
	{"intset_32", "dcd72eca6b84b6b0d53cc101bb942c93"},
	{"intset_64", "ffb6a23a556053311fbe9771dfe84b1c"},
	{"keys_with_expiry", "28096212e3df01ed8c4e97f4010ea424"},
	{"keys_with_mixed_expiry", "c7d02a3d9065abe01d212c6266a951c4"},
	{"linkedlist", "21ae916cd2aa9aaa79b242ac5c7c6510"},
	{"multiple_databases", "cbec71df9929b447e06d05218b7ec2ee"},
	{"rdb_version_5_with_checksum", "aca26c60a2d803610f32c80ef6d31e09"},
	{"regular_set", "2bbefdce0673c7038685212a77f3a926"},
	{"regular_sorted_set", "c3fd2dd7f5f90a97176b82900edc6f39"},
	{"sorted_set_as_ziplist", "0dba623f3c3ced434f95861cbc79a37f"},
	{"uncompressible_string_keys", "c2ee556c9740946c454a03b5f4a9a257"},
	{"ziplist_that_compresses_easily", "1fa6050f6eaa6ddcfcca67fd9ef161d8"},
	{"ziplist_that_doesnt_compress", "b236197a0e85bc6a6f0891ec8e79cc20"},
	{"ziplist_with_integers", "fff3650767ad19b377ed22aba8dc048b"},
	{"zipmap_that_compresses_easily", "7dccfe5b96da67f13aa6d42e29be04c6"},
	{"zipmap_that_doesnt_compress", "68744914aea876e9ecd012474f93dfc4"},
	{"zipmap_with_big_values", "cdb3480fb85e40a07518f75ee99d7705"},
}

func fixturePath(fixture string) string {
	return fmt.Sprintf("./fixtures/%s.rdb", fixture)
}

func TestFormatingFixturesUsingProtocol(t *testing.T) {
	var buf bytes.Buffer
	d := decoder.Protocol(&buf)

	for _, fixture := range fixtures {
		buf.Reset()
		format(fixturePath(fixture.name), d)

		md5sum := fmt.Sprintf("%x", md5.Sum(buf.Bytes()))
		if md5sum != fixture.md5Protocol {
			t.Errorf("%s: md5.Sum => %s want %s", fixture.name, md5sum, fixture.md5Protocol)
		}
	}
}
