package storage

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/GorunovAlx/shortening_long_url/internal/app/configs"
)

// FileStorage contains the path file.
type FileStorage struct {
	path string
}

// FileWriter contains a file for writing and bufio.Writer.
type FileWriter struct {
	file   *os.File
	writer *bufio.Writer
}

// FileScanner contains a file for scanning and bufio.Scanner.
type FileScanner struct {
	file    *os.File
	scanner *bufio.Scanner
}

// Returns a pointer to FileStorage with path file from config.
func NewInFileStorage() *FileStorage {
	return &FileStorage{
		path: configs.Cfg.FileStoragePath,
	}
}

// Returns a newly FileWriter.
func NewInFileWriter(st *FileStorage) (*FileWriter, error) {
	file, err := os.OpenFile(st.path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}

	return &FileWriter{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}

// Close file for FileWriter.
func (w *FileWriter) Close() error {
	return w.file.Close()
}

// Returns a newly FileScanner.
func NewInFileScanner(st *FileStorage) (*FileScanner, error) {
	file, err := os.OpenFile(st.path, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}

	return &FileScanner{
		file:    file,
		scanner: bufio.NewScanner(file),
	}, nil
}

// Close file for FileScanner.
func (w *FileScanner) Close() error {
	return w.file.Close()
}

// Writes a ShortURL to the file.
func (f *FileStorage) WriteShortURL(shortURL *ShortURL) error {
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

// Find and read shortened link and returns ShortURL.
func (f *FileStorage) ReadShortURL(shortLink string) (*ShortURL, error) {
	sc, err := NewInFileScanner(f)
	if err != nil {
		return nil, err
	}
	defer sc.Close()

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
