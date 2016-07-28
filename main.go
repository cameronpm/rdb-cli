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

var Usage = func() {
	fmt.Fprintf(os.Stderr, "Usage: %s [parameters]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\nParameters:\n")
	flag.PrintDefaults()
}

func availableFormats() []string {
	keys := make([]string, len(formats))

	i := 0
	for key := range formats {
		keys[i] = key
		i++
	}

	return keys
}

func main() {
	flag.Usage = Usage
	format := flag.String("format", "protocol", "Output format")
	rdbpath := flag.String("rdb", "", "Path to RDB file")
	flag.Parse()

	if len(*rdbpath) == 0 {
		fmt.Printf("Missing RDB file path\n\n")
		flag.Usage()
		os.Exit(1)
	}

	decoder, ok := formats[*format]
	if !ok {
		fmt.Printf("Invalid -format %s; Available formats: %s\n\n", *format, availableFormats())
		flag.Usage()
		os.Exit(1)
	}

	file, err := os.Open(*rdbpath)
	if err != nil {
		fmt.Printf("Fatal: %s\n", err)
		os.Exit(1)
	}

	err = rdb.Decode(file, decoder)
	if err != nil {
		fmt.Printf("Fatal: %s\n", err)
		os.Exit(1)
	}
}
