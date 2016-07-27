package main

import (
	"flag"
	"fmt"
	"github.com/brainly/rdb-cli/decoder"
	"github.com/cupcake/rdb"
	"os"
)

var formats = map[string]rdb.Decoder{
	"protocol": decoder.NewProtocolDecoder(os.Stdout),
	"diff":     &decoder.Diff{},
}

func main() {
	format := flag.String("format", "protocol", "Output format")
	rdbpath := flag.String("rdb", "", "Path to RDB file")
	flag.Parse()

	if len(*rdbpath) == 0 {
		fmt.Println("Missing RDB file path")
		flag.PrintDefaults()
		os.Exit(1)
	}

	decoder, ok := formats[*format]
	if !ok {
		fmt.Println("Invalid format ", *format)
		flag.PrintDefaults()
		os.Exit(1)
	}

	file, err := os.Open(*rdbpath)
	if err != nil {
		fmt.Printf("Fatal: %s", err)
		os.Exit(1)
	}

	err = rdb.Decode(file, decoder)
	if err != nil {
		fmt.Printf("Fatal: %s", err)
		os.Exit(1)
	}
}
