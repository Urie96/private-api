package main

import (
	// "fmt"

	"fmt"
	"io"

	cache "github.com/chenyahui/gin-cache"
	"github.com/gin-gonic/gin"
)

func main() {
	envs := GetEnvs()

	app := gin.Default()

	cacheFunc := cache.CacheByRequestURI(NewFileCacheStore(".cache"), 0)

	app.GET(
		"/tts",
		func(c *gin.Context) {
			p := TTSParams{
				VoiceType: 101016,
			}
			if c.BindQuery(&p) != nil {
				return
			}
			path, err := TencentTTSWithCache(p)
			if err != nil {
				c.AbortWithError(500, err)
				return
			}
			c.File(path)
		})
	app.GET(
		"/openai",
		cacheFunc,
		func(c *gin.Context) {
			p := OpenAIParams{
				Model:            "text-ada-001",
				Temperature:      0.9,
				MaxTokens:        150,
				TopP:             1,
				FrequencyPenalty: 0,
				PresencePenalty:  0.6,
				Stream:           true,
			}
			if c.BindQuery(&p) != nil {
				return
			}
			reader, err := OpenAIRequest(c.Request.Context(), p)
			if err != nil {
				c.AbortWithError(500, err)
				return
			}
			buf := make([]byte, 256)
			c.Stream(func(w io.Writer) bool {
				n, err := reader.Read(buf)
				if err != nil {
					if err != io.EOF {
						fmt.Println("read error:", err)
					}
					return false
				}
				w.Write(buf[:n])
				return true
			})
		},
	)

	err := app.Run(fmt.Sprintf(":%d", envs.Port))
	if err != nil {
		panic(err)
	}
}
