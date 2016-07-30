package main

import (
	"fmt"
	"github.com/brainly/rdb-cli/decoder"
	"github.com/cupcake/rdb"
	"github.com/jessevdk/go-flags"
	"log"
	"os"
)

var Version = "0.1-dev"

var formats = map[string]rdb.Decoder{
	"protocol": decoder.Protocol(os.Stdout),
	"diff":     decoder.Diff(),
}

type rdbOptions struct {
	Path string `required:"1" positional-arg-name:"RDBPATH" description:"Path to RDB file"`
}

var options struct {
	Version func()     `long:"version" description:"Print version and exit"`
	Format  string     `long:"format" choice:"diff" choice:"protocol" description:"Output format of RDB file"`
	Rdb     rdbOptions `positional-args:"1"`
}

func main() {
	options.Version = func() {
		fmt.Printf("rdb-cli v%s\n", Version)
		os.Exit(0)
	}

	if _, err := flags.Parse(&options); err != nil {
		os.Exit(1)
	}

	file, err := os.Open(options.Rdb.Path)
	if err != nil {
		log.Fatalf("Unable to open RDB file '%s': %s\n", options.Rdb.Path, err)
	}

	err = rdb.Decode(file, formats[options.Format])
	if err != nil {
		log.Fatalf("Failed to decode RDB file '%s': %s\n", options.Rdb.Path, err)
	}
}
