package main

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

func MD5(o any) string {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	// Encode (send) the value.
	err := enc.Encode(o)
	if err != nil {
		log.Fatal("gob encode error:", err)
	}

	sum := md5.Sum(b.Bytes())
	return hex.EncodeToString(sum[:])
}

func DecodeBase64(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

type Envs struct {
	TencentSDK struct {
		SecretKey    string `env:"TENCENT_SECRET_KEY,notEmpty"`
		SecretID     string `env:"TENCENT_SECRET_ID,notEmpty"`
		TTSProjectID int64  `env:"TENCENT_TTS_PROJECT_ID,notEmpty"`
	}
	Port   int64 `env:"PORT" envDefault:"8080"`
	OpenAI struct {
		SecretKey string `env:"OpenAI_SECRET_KEY,notEmpty"`
	}
}

var GetEnvs = func() func() Envs {
	godotenv.Load()

	envs := Envs{}
	if err := env.Parse(&envs); err != nil {
		panic(err)
	}
	return func() Envs {
		return envs
	}
}()

func P[T any](t T) *T {
	return &t
}

func NewFileCacheStore(cacheDir string) persist.CacheStore {
	os.MkdirAll(cacheDir, 0777)
	return FileCacheStore{strings.TrimRight(cacheDir, "/")}
}

type FileCacheStore struct {
	Dir string
}

func (c FileCacheStore) Get(key string, value interface{}) error {
	name := fmt.Sprintf("%s/%s", c.Dir, MD5(key))
	f, _ := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0644)
	defer f.Close()
	dec := gob.NewDecoder(f)
	if err := dec.Decode(value); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (c FileCacheStore) Set(key string, value interface{}, expire time.Duration) error {
	name := fmt.Sprintf("%s/%s", c.Dir, MD5(key))
	f, _ := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0644)
	defer f.Close()
	enc := gob.NewEncoder(f)
	if err := enc.Encode(value); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (FileCacheStore) Delete(key string) error {
	return nil
}
