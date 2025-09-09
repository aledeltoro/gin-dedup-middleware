package main

import (
	"log"
	"net/http"
	"time"

	"github.com/aledeltoro/gin-dedup-middleware/storage"
	"github.com/gin-gonic/gin"
)

func Deduplicate(cache storage.CacheStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		fullURL := c.Request.URL.String()
		param := c.Param("id")

		isDuplicateRequest, err := cache.IsSetMember(c, param, fullURL)
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

		err = cache.AddSet(c, param, 2*time.Minute, fullURL)
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
	router.Use(Deduplicate(redisCache))

	router.GET("/ping/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]any{
			"message": "pong",
			"id":      c.Param("id"),
		})
	})

	router.GET("/products/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]any{
			"message": "we have apples!",
			"id":      c.Param("id"),
		})
	})

	router.Run(":8080")
}
