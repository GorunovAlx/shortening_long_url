package storage

import (
	"errors"
	"strconv"
	"sync"
	//"github.com/caarlos0/env/v6"
)

type StorageConfig struct {
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:""`
}

var Cfg StorageConfig

// ShortURL struct contains a short link and initial link.
type ShortURL struct {
	ShortLink   string `json:"result,omitempty" valid:"-"`
	InitialLink string `json:"url,omitempty" valid:"-"`
}

// ShortURLRepo is an interface that contains two methods.
// GetInitialLink takes a short reference and returns the original.
// CreateShortURL takes an initial reference and returns a short.
type ShortURLRepo interface {
	GetInitialLink(shortLink string) (string, error)
	CreateShortURL(initialLink string) (string, error)
}

type ReadWriteShortURL interface {
	ReadShortURL(shortLink string) (*ShortURL, error)
	WriteShortURL(shortURL *ShortURL) error
}

// The ShortURLStorage contains data about the next short link,
// a repository with the type of map and mutex.
type ShortURLStorage struct {
	nextShortLink int
	storage       ReadWriteShortURL
	s             sync.RWMutex
}

// NewShortURLStorage returns a newly initialized ShortURLStorage object.
func NewShortURLStorage() (*ShortURLStorage, error) {
	/*
		err := env.Parse(&Cfg)
		if err != nil {
			log.Fatal(err)
		}
	*/
	if Cfg.FileStoragePath == "" {
		return &ShortURLStorage{
			nextShortLink: 1,
			storage:       NewInMemoryStorage(),
		}, nil
	} else {
		st, err := NewInFileStorage()

		if err != nil {
			return nil, errors.New("error occured in creating or opening file")
		}

		return &ShortURLStorage{
			nextShortLink: 1,
			storage:       st,
		}, nil
	}
}

// Get initial link by short link.
func (repo *ShortURLStorage) GetInitialLink(shortLink string) (string, error) {
	repo.s.RLock()
	defer repo.s.RUnlock()

	url, err := repo.storage.ReadShortURL(shortLink)
	if err != nil {
		return "", errors.New(err.Error())
	}

	return url.InitialLink, nil
}

// Create short link by initial link.
func (repo *ShortURLStorage) CreateShortURL(initialLink string) (string, error) {
	repo.s.Lock()
	defer repo.s.Unlock()

	sL := strconv.Itoa(repo.nextShortLink)
	shortURL := ShortURL{
		ShortLink:   sL,
		InitialLink: initialLink,
	}

	err := repo.storage.WriteShortURL(&shortURL)
	if err != nil {
		return "", errors.New(err.Error())
	}
	repo.nextShortLink += 1

	return shortURL.ShortLink, nil
}
