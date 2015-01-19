package main

import (
  "crypto/rand"
  "encoding/hex"
  "fmt"
  "github.com/garyburd/redigo/redis"
  "github.com/gorilla/mux"
  "github.com/soveran/redisurl"
  "net/http"
  "os"
)

func main() {
  route := mux.NewRouter()
  route.HandleFunc("/{code}", RedirectToLongURL).Methods("GET")
  route.HandleFunc("/", RenderForm).Methods("GET")
  route.HandleFunc("/", MakeShortURL).Methods("POST")

  http.Handle("/", route)

  fmt.Println("Listening...")
  err := http.ListenAndServe(":" + os.Getenv("PORT"), nil)
  if err != nil {
    panic(err)
  }
}

func MakeShortURL(res http.ResponseWriter, req *http.Request) {
  longUrl := req.FormValue("url")
  code, _ := randomHex(5)
  client, _ := redisurl.ConnectToURL(os.Getenv("REDISCLOUD_URL"))
  client.Do("SET", code, longUrl)
  client.Close()
  fmt.Fprintln(res, os.Getenv("HOST") + "/" + code)
}

func RedirectToLongURL(res http.ResponseWriter, req *http.Request) {
  code := req.URL.Path[1:]
  client, _ := redisurl.ConnectToURL(os.Getenv("REDISCLOUD_URL"))
  longUrl, _ := redis.String(client.Do("GET", code))
  client.Close()

  http.Redirect(res, req, longUrl, 302)
}

func RenderForm(res http.ResponseWriter, req *http.Request) {
  fmt.Fprintln(res, "<html><body><form method=\"post\"><input name=\"url\" type=\"text\"><input type=\"submit\"></form></body>")
}

func randomHex(n int) (string, error) {
  bytes := make([]byte, n)
  if _, err := rand.Read(bytes); err != nil {
    return "", err
  }
  return hex.EncodeToString(bytes), nil
}

