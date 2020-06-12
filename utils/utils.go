package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"sync"
)

const (
	// BufferSize is the size of the chan buffer
	BufferSize = 100
)

// Marshal is used to save the key
var Marshal = func(v interface{}) (io.Reader, error) {
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

// Unmarshal is a function that unmarshals the data from the
// reader into the specified value.
// By default, it uses the JSON unmarshaller.
var Unmarshal = func(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

type chunk struct {
	bufsize int
	offset  int64
}

// WriteCompressed creates the directory that stores the compressed file
func WriteCompressed(src, file string, key map[string]int) error {
	var filename string
	for _, v := range file {
		if v == '/' {
			fileslice := strings.Split(file, "/")
			filename = fileslice[len(fileslice)-1]
			break
		} else {
			filename = file
			continue
		}
	}

	path := file + ".fldr"

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		err = os.RemoveAll(path + "/")
		if err != nil {
			return err
		}
	}

	err := os.Mkdir(path, 0777)
	if err != nil {
		return err
	}
	keypath := path + "/" + ".keystore"
	f, err := os.Create(keypath)
	if err != nil {
		return err
	}
	defer f.Close()

	r, err := Marshal(key)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, r)
	if err != nil {
		return err
	}

	filepath := path + "/" + filename + ".cmp"
	_, err = os.Create(filepath)
	if err != nil {
		return err
	}

	err = WriteFile(src, filepath)
	if err != nil {
		return err
	}

	err = os.Remove(file)

	return err
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

// LoadCompressed loads the keystore
func LoadCompressed(file string, key interface{}) (string, error) {
	var filename string
	for _, v := range file {
		if v == '/' {
			fileslice := strings.Split(file, "/")
			filename = fileslice[len(fileslice)-1]
			break
		} else {
			filename = file
			continue
		}
	}
	path := file + ".fldr/.keystore"

	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	err = Unmarshal(f, key)
	if err != nil {
		return "", err
	}

	filepath := file + ".fldr/" + filename + ".cmp"
	s, err := LoadFile(filepath)
	if err != nil {
		return "", err
	}

	err = os.RemoveAll(file + ".fldr/")

	return s, err
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
