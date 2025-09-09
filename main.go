package main

import (
	"net/http"

	"github.com/aledeltoro/gin-dedup-middleware/storage"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var deduplicatedRequests = map[string]bool{}

func Deduplicate(redisCache storage.CacheStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		fullURL := c.Request.URL.String()

		_, ok := deduplicatedRequests[fullURL]
		if ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{
				"error":     "duplicate request",
				"full_path": fullURL,
			})

			return
		}

		deduplicatedRequests[fullURL] = true

		c.Next()
	}
}

func main() {
	redisCache := storage.NewRedisCache()

	router := gin.Default()
	router.Use(Deduplicate(redisCache))

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]any{
			"message":    "pong",
			"request_id": uuid.Must(uuid.NewV7()),
		})
	})

	router.Run(":8080")
}
