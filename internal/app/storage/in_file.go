package storage

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/GorunovAlx/shortening_long_url/internal/app/configs"
)

type InFileStorage struct {
	path string
}

type InFileWriter struct {
	file   *os.File
	writer *bufio.Writer
}

type InFileScanner struct {
	file    *os.File
	scanner *bufio.Scanner
}

func NewInFileStorage() (*InFileStorage, error) {
	return &InFileStorage{
		path: configs.Cfg.FileStoragePath,
	}, nil
}

func NewInFileWriter(st *InFileStorage) (*InFileWriter, error) {
	file, err := os.OpenFile(st.path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}

	return &InFileWriter{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}

func (w *InFileWriter) Close() error {
	return w.file.Close()
}

func NewInFileScanner(st *InFileStorage) (*InFileScanner, error) {
	file, err := os.OpenFile(st.path, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}

	return &InFileScanner{
		file:    file,
		scanner: bufio.NewScanner(file),
	}, nil
}

func (w *InFileScanner) Close() error {
	return w.file.Close()
}

func (f *InFileStorage) WriteShortURL(shortURL *ShortURL) error {
	data, err := json.Marshal(&shortURL)
	if err != nil {
		return err
	}

	wr, err := NewInFileWriter(f)
	if err != nil {
		return err
	}
	defer wr.Close()

	if _, err := wr.writer.Write(data); err != nil {
		return err
	}

	if err := wr.writer.WriteByte('\n'); err != nil {
		return err
	}

	return wr.writer.Flush()
}

func (f *InFileStorage) ReadShortURL(shortLink string) (*ShortURL, error) {
	sc, err := NewInFileScanner(f)
	if err != nil {
		return nil, err
	}
	defer sc.Close()

	/*
		if !sc.scanner.Scan() {
			return nil, sc.scanner.Err()
		}
	*/

	for sc.scanner.Scan() {
		data := sc.scanner.Bytes()
		shortURL := ShortURL{}
		err := json.Unmarshal(data, &shortURL)
		if err != nil {
			return nil, err
		}
		if shortURL.ShortLink == shortLink {
			return &shortURL, nil
		}
	}

	return &ShortURL{}, nil
}
