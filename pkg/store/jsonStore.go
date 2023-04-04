package store

import (
	"encoding/json"
	"os"
)

type Store interface {
	Read(data any) error
	Write(data any) error
	Check() error
}

type fileStore struct {
	Path string
}

func NewStore(fileName string) Store {
	return &fileStore{fileName}
}

func (fs *fileStore) Check() error {
	if _, err := os.ReadFile(fs.Path); err != nil {
		return err
	}
	
	return nil
}

func (fs *fileStore) Write(data any) error {
	dataFromFile, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return os.WriteFile(fs.Path, dataFromFile, 0644)
}

func (fs *fileStore) Read(data any) error {
	file, err := os.ReadFile(fs.Path)
	if err != nil {
		return err
	}

	return json.Unmarshal(file, &data)
}
