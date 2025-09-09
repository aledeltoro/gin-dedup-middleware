package dedup

import "github.com/gin-gonic/gin"

type DeduplicationOption int

const (
	WithParam DeduplicationOption = iota
	WithQuery
	WithHeader
)

type Config struct {
	Option DeduplicationOption
	Input  string
}

func NewDeduplicationKey(option DeduplicationOption, input string) *Config {
	return &Config{
		Option: option,
		Input:  input,
	}
}

func (d Config) Fetch(c *gin.Context) string {
	result := ""

	switch d.Option {
	case WithParam:
		result = c.Param(d.Input)
	case WithQuery:
		result = c.Query(d.Input)
	case WithHeader:
		result = c.GetHeader(d.Input)
	}

	return result
}
