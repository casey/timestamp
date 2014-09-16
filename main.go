package app

import "appengine"
import "fmt"
import "net/http"
import "regexp"
import "time"

import . "flotilla"

var path_re = regexp.MustCompile(`^/(?P<key>[a-zA-Z0-9_.-]*)$`)

func formatTime(t time.Time) string {
  return fmt.Sprintf("%f", float64(t.UnixNano())/1e9)
}

func init() {
  Handle("/").Get(get).Put(put)
}

func get(r *http.Request) {
  c := appengine.NewContext(r)
  match := Components(path_re, r.URL.Path)
  Ensure(match != nil, http.StatusForbidden)
  pointer, e := getTimestamp(c, match["key"])
  Check(e)
  Ensure(pointer != nil, http.StatusNotFound)
  Body(http.StatusOK, formatTime(*pointer), "text/plain; charset=utf-8")
}

func put(r *http.Request) {
  c := appengine.NewContext(r)
  match := Components(path_re, r.URL.Path)
  Ensure(match != nil, http.StatusForbidden)
  key := match["key"]

  time, e := getTimestamp(c, key)
  Check(e)

  if time != nil {
    Body(http.StatusOK, formatTime(*time), "text/plain; charset=utf-8")
  }

  time, e = putTimestamp(c, key)
  Check(e)
  Body(http.StatusCreated, formatTime(*time), "text/plain; charset=utf-8")
}
