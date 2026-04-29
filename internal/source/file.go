package source

import (
	"encoding/binary"
	"os"
)

type FileSource struct {
	file *os.File
}

func ReadFile(path string) (*FileSource, error) {
	openedFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return &FileSource{file: openedFile}, nil
}

func (f *FileSource) Read(buffer []int8) error {
	err := binary.Read(f.file, binary.LittleEndian, &buffer)
	return err
}

func (f *FileSource) Close() error {
	return f.file.Close()
}
