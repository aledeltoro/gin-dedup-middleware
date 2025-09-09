package main

import (
	"log"
	"net/http"
	"time"

	"github.com/aledeltoro/gin-dedup-middleware/dedup"
	"github.com/aledeltoro/gin-dedup-middleware/storage"
	"github.com/gin-gonic/gin"
)

func Deduplicate(cache storage.CacheStorage, dedupConfig *dedup.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		fullURL := c.Request.URL.String()
		dedupKey := dedupConfig.Fetch(c)

		isDuplicateRequest, err := cache.IsSetMember(c, dedupKey, fullURL)
		if err != nil {
			log.Printf("failed performing SISMEMBER command: %s", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{
				"error": "internal server error",
			})
			return
		}

		if isDuplicateRequest {
			c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{
				"error":     "duplicate request",
				"full_path": fullURL,
			})

			return
		}

		err = cache.AddSet(c, dedupKey, 2*time.Minute, fullURL)
		if err != nil {
			log.Printf("failed performing SADD command: %s", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{
				"error": "internal server error",
			})
		}

		c.Next()
	}
}

func main() {
	redisCache := storage.NewRedisCache()

	router := gin.Default()

	router.GET(
		"/ping/:id",
		Deduplicate(redisCache, dedup.NewDeduplicationKey(dedup.WithParam, "id")),
		func(c *gin.Context) {
			c.JSON(http.StatusOK, map[string]any{
				"message": "pong",
				"id":      c.Param("id"),
			})
		})

	router.GET(
		"/products/:id",
		Deduplicate(redisCache, dedup.NewDeduplicationKey(dedup.WithParam, "id")),
		func(c *gin.Context) {
			c.JSON(http.StatusOK, map[string]any{
				"message": "we have apples!",
				"id":      c.Param("id"),
			})
		})

	router.Run(":8080")
}
