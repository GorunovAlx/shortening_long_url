package storage

import (
	"errors"
	"sync"

	"github.com/GorunovAlx/shortening_long_url/internal/app/configs"
	gen "github.com/GorunovAlx/shortening_long_url/internal/app/generators"
)

// ShortURL struct contains a InitialLink - initial link
// and its shortened link(ShortLink).
type ShortURL struct {
	InitialLink string `json:"url,omitempty" valid:"-"`
	ShortLink   string `json:"result,omitempty" valid:"-"`
	UserID      uint32 `json:"user_id,omitempty"`
}

type ShortURLByUser struct {
	ShortLink   string `json:"short_url,omitempty" valid:"-"`
	InitialLink string `json:"original_url,omitempty" valid:"-"`
}

// ShortURLRepo contains:
// GetInitialLink takes a short link and returns the initial link;
// CreateShortURL takes an initial link and returns a shortened.
type ShortURLRepo interface {
	GetInitialLink(shortLink string) (string, error)
	CreateShortURL(shortURL *ShortURL) (string, error)
	GetAllShortURLUser(id uint32) ([]ShortURLByUser, error)
}

// RWShortURL contains:
// GetInitialLink takes a short link and returns the initial link from storage;
// WriteShortURL takes the ShortURL struct and writes it into the storage.
type StorageOperations interface {
	GetInitialLink(shortLink string) (string, error)
	WriteShortURL(shortURL *ShortURL) error
	GetAllShortURLByUser(userID uint32) ([]ShortURLByUser, error)
}

// The ShortURLStorage contains storage that implements
// the interface RWShortURL and RWMutex.
type ShortURLStorage struct {
	storage StorageOperations
	s       sync.RWMutex
}

// The function returns a pointer to the ShortURLStorage structure,
// where the storage is initialized either by file storage
// if the file path is not empty in the config, or by in memory storage.
func NewStorage() *ShortURLStorage {
	if configs.Cfg.FileStoragePath != "" {
		return &ShortURLStorage{
			storage: NewInFileStorage(),
		}
	}

	return &ShortURLStorage{
		storage: NewInMemoryStorage(),
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
func (repo *ShortURLStorage) CreateShortURL(shortURL *ShortURL) (string, error) {
	repo.s.Lock()
	defer repo.s.Unlock()

	shortenedURL, e := gen.GenerateShortLink(shortURL.InitialLink)
	if e != nil {
		return "", errors.New(e.Error())
	}

	shortURL.ShortLink = shortenedURL
	err := repo.storage.WriteShortURL(shortURL)
	if err != nil {
		return "", errors.New(err.Error())
	}

	return shortURL.ShortLink, nil
}

func (repo *ShortURLStorage) GetAllShortURLUser(id uint32) ([]ShortURLByUser, error) {
	repo.s.RLock()
	defer repo.s.RUnlock()

	result, err := repo.storage.GetAllShortURLByUser(id)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return result, nil
}
