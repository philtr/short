package main

import (
  "crypto/rand"
  "encoding/hex"
  "fmt"
  "github.com/fzzy/radix/redis"
  "github.com/gorilla/mux"
  "net/http"
  "os"
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
  client, _ := redis.Dial("tcp", os.Getenv("REDISCLOUD_URL"))
  longUrl, _ := client.Cmd("GET", code).Str()
  http.Redirect(res, req, longUrl, 302)
}

func short(res http.ResponseWriter, req *http.Request) {
  longUrl := req.FormValue("url")
  fmt.Fprintln(res, makeShort(longUrl))
}

func makeShort(longUrl string) (string) {
  code, _ := randomHex(5)
  client, _ := redis.Dial("tcp", os.Getenv("REDISCLOUD_URL"))
  client.Cmd("SET", code, longUrl)
  return os.Getenv("HOST") + "/" + code
}

func randomHex(n int) (string, error) {
  bytes := make([]byte, n)
  if _, err := rand.Read(bytes); err != nil {
    return "", err
  }
  return hex.EncodeToString(bytes), nil
}


