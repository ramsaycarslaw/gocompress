# Golang Compression

This package allows for the compression of text from golang with comparable size reduction to gzip

## Instillation

The package can be installed with the following

```
go get github.com/ramsaycarslaw/compression
```

## Usage

To create a simple command line tool to compress then decompress a file

```
package main

import (
	"cmp/compression"
	"log"
	"os"
)

func main() {
	var c = compression.NewCompressor(os.Args[1])

	err := c.Compress()
	if err != nil {
		log.Fatal(err)
	}

	err = c.Decompress()
	if err != nil {
		log.Fatal(err)
	}
}
```

And then run it with

```
./compression file.txt
```

This will create a new file called file.txt.cmp which is the compressed file
