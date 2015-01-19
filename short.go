package main

import (
  "crypto/rand"
  "encoding/hex"
  "fmt"
  "github.com/gorilla/mux"
  "net/http"
  "os"
  "github.com/soveran/redisurl"
  "github.com/garyburd/redigo/redis"
)

func main() {
  r := mux.NewRouter()
  r.HandleFunc("/", short).Methods("POST")
  r.HandleFunc("/{code}", redir)

  http.Handle("/", r)

  fmt.Println("Listening...")
  err := http.ListenAndServe(":" + os.Getenv("PORT"), nil)
  if err != nil {
    panic(err)
  }
}

func redir(res http.ResponseWriter, req *http.Request) {
  code := req.URL.Path[1:]

  client, _ := redisurl.ConnectToURL(os.Getenv("REDISCLOUD_URL"))
  longUrl, _ := redis.String(client.Do("GET", code))
  client.Close()

  http.Redirect(res, req, longUrl, 302)
}

func short(res http.ResponseWriter, req *http.Request) {
  longUrl := req.FormValue("url")
  fmt.Fprintln(res, makeShort(longUrl))
}

func makeShort(longUrl string) (string) {
  code, _ := randomHex(5)
  client, _ := redisurl.ConnectToURL(os.Getenv("REDISCLOUD_URL"))
  client.Do("SET", code, longUrl)
  client.Close()
  return os.Getenv("HOST") + "/" + code
}

func randomHex(n int) (string, error) {
  bytes := make([]byte, n)
  if _, err := rand.Read(bytes); err != nil {
    return "", err
  }
  return hex.EncodeToString(bytes), nil
}


