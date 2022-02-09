package storage

import (
	"errors"
	"sync"
)

type ShortURL struct {
	ShortLink   int
	InitialLink string
}

type ShortURLRepo interface {
	GetInitialLink(shortLink int) (string, error)
	CreateShortURL(initialLink string) (int, error)
}

type ShortURLStorage struct {
	nextShortLink int
	storage       map[int]ShortURL
	s             sync.RWMutex
}

func NewShortURLStorage() *ShortURLStorage {
	return &ShortURLStorage{
		nextShortLink: 1,
		storage:       make(map[int]ShortURL),
	}
}

func (repo *ShortURLStorage) GetInitialLink(shortLink int) (string, error) {
	repo.s.RLock()
	defer repo.s.RUnlock()

	url, ok := repo.storage[shortLink]
	if !ok {
		return "", errors.New("the url with this value does not exist")
	}

	return url.InitialLink, nil
}

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
