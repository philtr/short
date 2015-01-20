package main

import (
  "crypto/rand"
  "encoding/hex"
  "github.com/garyburd/redigo/redis"
  "github.com/soveran/redisurl"
  "os"
)

type Link struct {
  longUrl string
  code    string
}

func (l *Link) shortUrl() string {
  return os.Getenv("SHORTURL_HOST") + "/" + l.code
}

func (l *Link) save() {
  if l.code == "" {
    code := FindByLongUrl(l.longUrl).code
    if code == "" {
      l.code, _ = randomHex(2)
      go writeLink(l)
    }
  } else {
    l.longUrl = FindByCode(l.code).longUrl
  }
}

func FindByCode(code string) Link {
  client, _ := redisurl.ConnectToURL(os.Getenv("REDISCLOUD_URL"))
  longUrl, _ := redis.String(client.Do("GET", code))
  client.Close()

  return Link{code: code, longUrl: longUrl}
}

func FindByLongUrl(longUrl string) Link {
  client, _ := redisurl.ConnectToURL(os.Getenv("REDISCLOUD_URL"))
  code, _ := redis.String(client.Do("GET", longUrl))
  client.Close()

  return Link{code: code, longUrl: longUrl}
}

func randomHex(n int) (string, error) {
  bytes := make([]byte, n)
  if _, err := rand.Read(bytes); err != nil {
    return "", err
  }
  return hex.EncodeToString(bytes), nil
}

func writeLink(l *Link) {
  client, _ := redisurl.ConnectToURL(os.Getenv("REDISCLOUD_URL"))
  client.Do("SET", l.code, l.longUrl) // find longUrl by code
  client.Do("SET", l.longUrl, l.code) // find code by longUrl
  client.Close()
}


