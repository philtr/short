package main

import (
  "crypto/rand"
  "encoding/hex"
  "github.com/garyburd/redigo/redis"
  "github.com/go-martini/martini"
  "github.com/martini-contrib/render"
  "github.com/soveran/redisurl"
  "net/http"
  "os"
)

func main() {
  m := martini.Classic()
  m.Use(render.Renderer())

  m.Get("/", RenderForm)
  m.Get("/:code", RedirectToLongURL)
  m.Post("/", MakeShortURL)

  m.Run()
}

func MakeShortURL(req *http.Request) (string) {
  longUrl := req.FormValue("url")
  code := saveLink(longUrl)
  shortUrl := os.Getenv("SHORTURL_HOST") + "/" + code

  return shortUrl
}

func RedirectToLongURL(params martini.Params, r render.Render) {
  r.Redirect(findLink(params["code"]), 302)
}

func RenderForm(r render.Render) {
  r.HTML(200, "form", "")
}

func findCode(longUrl string) (string) {
  client, _ := redisurl.ConnectToURL(os.Getenv("REDISCLOUD_URL"))
  code, _ := redis.String(client.Do("GET", longUrl))
  client.Close()

  return code
}


func findLink(code string) (string) {
  client, _ := redisurl.ConnectToURL(os.Getenv("REDISCLOUD_URL"))
  longUrl, _ := redis.String(client.Do("GET", code))
  client.Close()

  return longUrl
}

func randomHex(n int) (string, error) {
  bytes := make([]byte, n)
  if _, err := rand.Read(bytes); err != nil {
    return "", err
  }
  return hex.EncodeToString(bytes), nil
}

func saveLink(longUrl string) (string){
  code := findCode(longUrl)

  if code == "" {
    code, _ := randomHex(2)
    go writeLink(code, longUrl)
    return code
  }

  return code
}

func writeLink(code string, longUrl string) {
  client, _ := redisurl.ConnectToURL(os.Getenv("REDISCLOUD_URL"))
  client.Do("SET", code, longUrl)
  client.Do("SET", longUrl, code)
  client.Close()
}
