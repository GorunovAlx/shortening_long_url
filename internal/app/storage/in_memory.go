package storage

import (
	"errors"
)

// InMemoryStorage contains storage map[string]ShortURL.
type InMemoryStorage struct {
	storage map[string]string
}

// Returns a pointer to InMemoryStorage.
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		storage: make(map[string]string),
	}
}

// Find and read shortened link and returns ShortURL.
func (m *InMemoryStorage) GetInitialLink(shortLink string) (string, error) {
	initialLink, ok := m.storage[shortLink]
	if !ok {
		return "", errors.New("URL with this value does not exist")
	}
	return initialLink, nil
}

// Writes a ShortURL to the in memory storage.
func (m *InMemoryStorage) WriteShortURL(shortURL *ShortURL) error {
	for _, existing := range m.storage {
		if shortURL.InitialLink == existing {
			return errors.New("URL with same location already exists")
		}
	}
	m.storage[shortURL.ShortLink] = shortURL.InitialLink
	return nil
}
