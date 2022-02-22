package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	valid "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"

	"github.com/GorunovAlx/shortening_long_url/internal/app/storage"
)

func RegisterRoutes(repo storage.ShortURLRepo) *gin.Engine {
	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/:shortURL", GetInitialLinkHandler(repo))
	router.POST("/", CreateShortURLHandler(repo))
	router.POST("/api/shorten", CreateShortURLJSONHandler(repo))

	return router
}

func CreateShortURLJSONHandler(urlStorage storage.ShortURLRepo) func(c *gin.Context) {
	return func(c *gin.Context) {
		var url storage.ShortURL
		if err := json.NewDecoder(c.Request.Body).Decode(&url); err != nil {
			c.String(400, err.Error())
		}

		isURL := valid.IsURL(url.InitialLink)
		if !isURL {
			c.String(400, "Incorrect link")
		}

		shortURL, err := urlStorage.CreateShortURL(url.InitialLink)
		if err != nil {
			c.String(400, err.Error())
			return
		}

		res := storage.ShortURL{
			ShortLink: shortURL,
		}

		c.JSON(http.StatusCreated, res)
	}
}

func CreateShortURLHandler(urlStorage storage.ShortURLRepo) func(c *gin.Context) {
	return func(c *gin.Context) {
		b, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.String(400, err.Error())
		}
		if len(b) == 0 {
			c.String(400, "Incorrect request")
		}

		shortURL, err := urlStorage.CreateShortURL(string(b))
		if err != nil {
			c.String(400, err.Error())
		}

		c.String(http.StatusCreated, "http://localhost:8080/"+shortURL)
	}
}

func GetInitialLinkHandler(urlStorage storage.ShortURLRepo) func(c *gin.Context) {
	return func(c *gin.Context) {
		shortURL := c.Param("shortURL")
		if shortURL == "" {
			c.String(400, "short url was not sent")
		}

		link, err := urlStorage.GetInitialLink(shortURL)
		if err != nil {
			c.String(400, err.Error())
		}
		c.Header("Location", link)
		c.Status(http.StatusTemporaryRedirect)
	}
}
