package storage

import (
	"errors"
)

// InMemoryStorage contains storage map[string]ShortURL.
type InMemoryStorage struct {
	storage map[string]ShortURL
}

// Returns a pointer to InMemoryStorage.
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		storage: make(map[string]ShortURL),
	}
}

// Find and read shortened link and returns ShortURL.
func (m *InMemoryStorage) ReadShortURL(shortLink string) (*ShortURL, error) {
	sU, ok := m.storage[shortLink]
	if !ok {
		return nil, errors.New("the url with this value does not exist")
	}
	return &sU, nil
}

// Writes a ShortURL to the in memory storage.
func (m *InMemoryStorage) WriteShortURL(shortURL *ShortURL) error {
	for _, existing := range m.storage {
		if shortURL.InitialLink == existing.InitialLink {
			return errors.New("URL with same location already exists")
		}
	}
	m.storage[shortURL.ShortLink] = *shortURL
	return nil
}
