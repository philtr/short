package main

import (
  "github.com/go-martini/martini"
  "github.com/martini-contrib/render"
  "net/http"
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
  l := new(Link)
  l.longUrl = req.FormValue("url")
  l.save()

  return l.shortUrl()
}

func RedirectToLongURL(params martini.Params, r render.Render) {
  l := FindByCode(params["code"])
  r.Redirect(l.longUrl, 302)
}

func RenderForm(r render.Render) {
  r.HTML(200, "form", "")
}

