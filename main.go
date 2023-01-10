package main

import (
	// "fmt"

	"fmt"

	cache "github.com/chenyahui/gin-cache"
	"github.com/gin-gonic/gin"
)

func main() {
	envs := GetEnvs()

	app := gin.Default()

	cacheFunc := cache.CacheByRequestURI(NewFileCacheStore(".cache"), 0)

	app.GET(
		"/tts",
		cacheFunc,
		func(c *gin.Context) {
			p := TTSParams{
				VoiceType: 101016,
			}
			if c.BindQuery(&p) != nil {
				return
			}
			b, err := TencentTTS(p)
			if err != nil {
				c.AbortWithError(500, err)
				return
			}
			c.Data(200, "audio/x-wav", b)
		})

	err := app.Run(fmt.Sprintf(":%d", envs.Port))
	if err != nil {
		panic(err)
	}
}
