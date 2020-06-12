package utils

import (
	"io"
	"os"
	"sync"
)

const (
	// BufferSize is the size of the chan buffer
	BufferSize = 100
)

type chunk struct {
	bufsize int
	offset  int64
}

// WriteFile converts a string to bytes and writes it to a file
func WriteFile(src, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	b := []byte(src)
	_, err = f.Write(b)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	return nil
}

// LoadFile concurrently loads a file into a string
func LoadFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}

	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		return "", err
	}

	filesize := int(fileinfo.Size())

	// number of go routines
	concurrency := filesize / BufferSize

	chunksizes := make([]chunk, concurrency)

	for i := 0; i < concurrency; i++ {
		chunksizes[i].bufsize = BufferSize
		chunksizes[i].offset = int64(BufferSize * i)
	}

	// extra go routione in the case of remainder
	if remainder := filesize % BufferSize; remainder != 0 {
		c := chunk{bufsize: remainder, offset: int64(concurrency * BufferSize)}
		concurrency++
		chunksizes = append(chunksizes, c)
	}

	var out = make([][]byte, BufferSize*concurrency)
	var wg sync.WaitGroup
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func(chuncksizes []chunk, i int) {
			defer wg.Done()

			chunk := chuncksizes[i]
			buffer := make([]byte, chunk.bufsize)
			_, err := file.ReadAt(buffer, chunk.offset)

			if err != nil && err != io.EOF {
				return
			}

			//fmt.Println(string(buffer))
			out[i] = buffer
		}(chunksizes, i)
	}
	wg.Wait()
	var result string
	for _, v := range out {
		if v == nil {
			break
		}
		result += string(v)
	}
	return result, nil
}
