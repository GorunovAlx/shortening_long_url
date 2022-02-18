package storage

import (
	"errors"
)

type InMemory struct {
	st map[string]ShortURL
}

type InMemoryStorage struct {
	storage InMemory
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		storage: InMemory{
			st: make(map[string]ShortURL),
		},
	}
}

func (m *InMemoryStorage) ReadShortURL(shortLink string) (*ShortURL, error) {
	sU, ok := m.storage.st[shortLink]
	if !ok {
		return nil, errors.New("the url with this value does not exist")
	}
	return &sU, nil
}

func (m *InMemoryStorage) WriteShortURL(shortURL *ShortURL) error {
	for _, existing := range m.storage.st {
		if shortURL.InitialLink == existing.InitialLink {
			return errors.New("URL with same location already exists")
		}
	}
	m.storage.st[shortURL.ShortLink] = *shortURL
	return nil
}
