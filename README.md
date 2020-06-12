# Gol Compression

This package allows for the compression of text from golang with comparable size reduction to gzip

## Instillation

The package can be installed with the following

```
go get github.com/ramsaycarslaw/gocompress
```

## Usage

To create a simple command line tool to compress then decompress a file

```
package main

import (
	"log"
	"os"

	github.com/ramsaycarslaw/gocompress
)

func main() {
	var c = gocompress.NewCompressor(os.Args[1])

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
./compression path/to/file.txt
```

This will create a folder in the path called file.txt.fldr, inside of this is the compressed file and the keystore for the dictionary compression
