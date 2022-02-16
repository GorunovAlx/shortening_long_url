package storage

import (
	"errors"
	"sync"
)

// ShortURL struct contains a short link and initial link.
type ShortURL struct {
	ShortLink   int
	InitialLink string
}

// ShortURLRepo is an interface that contains two methods.
// GetInitialLink takes a short reference and returns the original.
// CreateShortURL takes an initial reference and returns a short.
type ShortURLRepo interface {
	GetInitialLink(shortLink int) (string, error)
	CreateShortURL(initialLink string) (int, error)
}

// The ShortURLStorage contains data about the next short link,
// a repository with the type of map and mutex.
type ShortURLStorage struct {
	nextShortLink int
	storage       map[int]ShortURL
	s             sync.RWMutex
}

// NewShortURLStorage returns a newly initialized ShortURLStorage object.
func NewShortURLStorage() *ShortURLStorage {
	return &ShortURLStorage{
		nextShortLink: 1,
		storage:       make(map[int]ShortURL),
	}
}

// Get initial link by short link.
func (repo *ShortURLStorage) GetInitialLink(shortLink int) (string, error) {
	repo.s.RLock()
	defer repo.s.RUnlock()

	url, ok := repo.storage[shortLink]
	if !ok {
		return "", errors.New("the url with this value does not exist")
	}

	return url.InitialLink, nil
}

// Create short link by initial link.
func (repo *ShortURLStorage) CreateShortURL(initialLink string) (int, error) {
	repo.s.Lock()
	defer repo.s.Unlock()

	for _, existing := range repo.storage {
		if initialLink == existing.InitialLink {
			return -1, errors.New("URL with same location already exists")
		}
	}

	shortURL := ShortURL{
		ShortLink:   repo.nextShortLink,
		InitialLink: initialLink,
	}

	repo.storage[shortURL.ShortLink] = shortURL
	repo.nextShortLink += 1

	return shortURL.ShortLink, nil
}
