package storage

import (
	"errors"
	"strconv"
	"sync"

	"github.com/GorunovAlx/shortening_long_url/internal/app/configs"
)

// ShortURL struct contains a InitialLink - initial link
// and its shortened link(ShortLink).
type ShortURL struct {
	InitialLink string `json:"url,omitempty" valid:"-"`
	ShortLink   string `json:"result,omitempty" valid:"-"`
}

// ShortURLRepo contains:
// GetInitialLink takes a short link and returns the initial link;
// CreateShortURL takes an initial link and returns a shortened.
type ShortURLRepo interface {
	GetInitialLink(shortLink string) (string, error)
	CreateShortURL(initialLink string) (string, error)
}

// RWShortURL contains:
// GetInitialLink takes a short link and returns the initial link from storage;
// WriteShortURL takes the ShortURL struct and writes it into the storage.
type RWShortURL interface {
	GetInitialLink(shortLink string) (string, error)
	WriteShortURL(shortURL *ShortURL) error
}

// The ShortURLStorage contains next short link,
// storage that implements the interface RWShortURL and RWMutex.
type ShortURLStorage struct {
	nextShortLink int
	storage       RWShortURL
	s             sync.RWMutex
}

// The function returns a pointer to the ShortURLStorage structure,
// where the storage is initialized either by file storage
// if the file path is not empty in the config, or by in memory storage.
func NewStorage() *ShortURLStorage {
	if configs.Cfg.FileStoragePath != "" {
		return &ShortURLStorage{
			nextShortLink: 1,
			storage:       NewInFileStorage(),
		}
	}

	return &ShortURLStorage{
		nextShortLink: 1,
		storage:       NewInMemoryStorage(),
	}
}

// Get the initial link by shortened link or an error.
func (repo *ShortURLStorage) GetInitialLink(shortLink string) (string, error) {
	repo.s.RLock()
	defer repo.s.RUnlock()

	url, err := repo.storage.GetInitialLink(shortLink)
	if err != nil {
		return "", errors.New(err.Error())
	}

	return url, nil
}

// Create shortened link by initial link.
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
