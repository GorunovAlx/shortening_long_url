package storage

import (
	"errors"
	"log"
	"sync"

	"github.com/GorunovAlx/shortening_long_url/internal/app/configs"
	gen "github.com/GorunovAlx/shortening_long_url/internal/app/generators"
	"github.com/GorunovAlx/shortening_long_url/internal/app/utils"
)

// ShortURL struct contains a InitialLink - initial link
// and its shortened link(ShortLink).
type ShortURL struct {
	InitialLink string `json:"url,omitempty" valid:"-"`
	ShortLink   string `json:"result,omitempty" valid:"-"`
	UserID      uint32 `json:"user_id,omitempty"`
}

type ShortURLByUser struct {
	ShortLink     string `json:"short_url,omitempty" valid:"-"`
	InitialLink   string `json:"original_url,omitempty" valid:"-"`
	CorrelationID string `json:"correlation_id,omitempty"`
}

// ShortURLRepo contains:
// GetInitialLink takes a short link and returns the initial link;
// CreateShortURL takes an initial link and returns a shortened.
type ShortURLRepo interface {
	GetInitialLink(shortLink string) (string, error)
	CreateShortURL(shortURL *ShortURL) (string, error)
	CreateListShortURL(links []ShortURLByUser) ([]ShortURLByUser, error)
	GetAllShortURLUser(id uint32) ([]ShortURLByUser, error)
	PingDB() error
}

// RWShortURL contains:
// GetInitialLink takes a short link and returns the initial link from storage;
// WriteShortURL takes the ShortURL struct and writes it into the storage.
type StorageOperations interface {
	GetInitialLink(shortLink string) (string, error)
	WriteShortURL(shortURL *ShortURL) error
	WriteListShortURL(links []ShortURLByUser) error
	GetAllShortURLByUser(userID uint32) ([]ShortURLByUser, error)
	PingDB() error
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
	if configs.Cfg.DatabaseDSN != "" {
		st, err := NewDBStorage()
		if err != nil {
			log.Println(err)
		} else {
			return &ShortURLStorage{
				storage: st,
			}
		}
	}

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
		return "", err
	}

	return url, nil
}

// Create shortened link by initial link.
func (repo *ShortURLStorage) CreateShortURL(shortURL *ShortURL) (string, error) {
	repo.s.Lock()
	defer repo.s.Unlock()

	shortenedURL, err := gen.GenerateShortLink(shortURL.InitialLink, shortURL.UserID)
	if err != nil {
		return "", err
	}
	shortURL.ShortLink = shortenedURL

	err = repo.storage.WriteShortURL(shortURL)

	if errors.Is(err, utils.ErrUniqueLink) {
		return shortURL.ShortLink, utils.ErrUniqueLink
	} else if err != nil {
		return "", err
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

func (repo *ShortURLStorage) PingDB() error {
	err := repo.storage.PingDB()

	if err != nil {
		return err
	}

	return nil
}

func (repo *ShortURLStorage) CreateListShortURL(links []ShortURLByUser) ([]ShortURLByUser, error) {
	repo.s.Lock()
	defer repo.s.Unlock()

	var shortenedLinks []ShortURLByUser

	for i, link := range links {
		shortenedURL, err := gen.GenerateShortLink(link.InitialLink, 0)
		if err != nil {
			return nil, err
		}
		shortenedURL = configs.Cfg.BaseURL + "/" + shortenedURL

		links[i].setShortLink(shortenedURL)
	}

	for _, link := range links {
		shortened := ShortURLByUser{
			ShortLink:     link.ShortLink,
			CorrelationID: link.CorrelationID,
		}

		shortenedLinks = append(shortenedLinks, shortened)
	}

	err := repo.storage.WriteListShortURL(links)
	if err != nil {
		return nil, err
	}

	return shortenedLinks, nil
}

func (s *ShortURLByUser) setShortLink(value string) {
	(*s).ShortLink = value
}
