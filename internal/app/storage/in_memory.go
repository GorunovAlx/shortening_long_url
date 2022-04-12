package storage

import (
	"errors"

	"github.com/GorunovAlx/shortening_long_url/internal/app/configs"
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
func (m *InMemoryStorage) GetInitialLink(shortLink string) (string, error) {
	sh, ok := m.storage[shortLink]
	if !ok {
		return "", errors.New("URL with this value does not exist")
	}
	return sh.InitialLink, nil
}

// Writes a ShortURL to the in memory storage.
func (m *InMemoryStorage) WriteShortURL(shortURL *ShortURL) error {
	for _, existing := range m.storage {
		if shortURL.InitialLink == existing.InitialLink {
			return nil
		}
	}
	m.storage[shortURL.ShortLink] = *shortURL
	return nil
}

func (m *InMemoryStorage) GetAllShortURLByUser(userID uint32) ([]ShortURLByUser, error) {
	var result []ShortURLByUser
	for _, shortURL := range m.storage {
		if shortURL.UserID == userID {
			byUser := ShortURLByUser{
				InitialLink: shortURL.InitialLink,
				ShortLink:   configs.Cfg.BaseURL + "/" + shortURL.ShortLink,
			}
			result = append(result, byUser)
		}
	}

	return result, nil
}

func (m *InMemoryStorage) PingDB() error {
	return errors.New("this type of storage does not support the ping operation")
}

func (m *InMemoryStorage) WriteListShortURL(links []ShortURLByUser) error {
	for _, link := range links {
		var url ShortURL
		url.InitialLink = link.InitialLink
		url.ShortLink = link.ShortLink
		m.storage[url.ShortLink] = url
	}

	return nil
}

func (m *InMemoryStorage) DeleteShortURLByUser(link string, id uint32) error {
	return nil
}
