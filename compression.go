// Copyright (C) 2020 Ramsay Carslaw

package gocompress

import (
	"errors"

	"github.com/ramsaycarslaw/gocompress/dict"
	"github.com/ramsaycarslaw/gocompress/utils"
	"github.com/ramsaycarslaw/gocompress/whitespace"
)

// Compressor is the main struct used to perform compression
type Compressor struct {
	DictKey  map[string]int
	Filename string
}

// NewCompressor is a initiliser for a compressor
func NewCompressor(filename string) *Compressor {
	cmp := &Compressor{
		Filename: filename,
	}
	return cmp
}

// Compress compresswes a file
func (c *Compressor) Compress() error {
	out, err := utils.LoadFile(c.Filename)
	if err != nil {
		return errors.New("Error reading file: " + c.Filename + " with error: " + err.Error())
	}

	result := whitespace.RemoveTrailingWhitespace(out)
	result = whitespace.ReplaceSpaces(result)
	result, m := dict.DictionaryCompression(result)
	c.DictKey = m

	err = utils.WriteCompressed(result, c.Filename, c.DictKey)
	if err != nil {
		return errors.New("Error writing to new file with error: " + err.Error())
	}

	return nil
}

// Decompress undoes the compression
func (c *Compressor) Decompress() error {
	s, err := utils.LoadCompressed(c.Filename, &c.DictKey)
	if err != nil {
		return errors.New("Error reading file: " + c.Filename + ".cmp" + " with error: " + err.Error())
	}

	result, err := dict.DictionaryDecompression(s, c.DictKey)
	if err != nil {
		return errors.New("Error reversing dict compression: " + err.Error())
	}

	err = utils.WriteFile(result, c.Filename)
	if err != nil {
		return errors.New("Error writing to new file with error: " + err.Error())
	}

	return nil
}
