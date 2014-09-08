package timestamp

import "appengine"
import "appengine/datastore"
import "crypto/sha256"
import "encoding/hex"
import "fmt"
import "net/http"
import "regexp"
import "time"

var path_re = regexp.MustCompile(`^/([a-zA-Z0-9.-_]*)$`)

type Entity struct {
  Timestamp time.Time
}

func init() {
  http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)

  body := ""

  defer func() {
    relic := recover()
    c.Infof("relic: %v", relic)
    code, ok := relic.(int)
    if ok {
      if body == "" {
        body = http.StatusText(code);
      }
      w.Header().Set("Warranty", `THIS IS PROVIDED "AS IS" WITHOUT WARRANTY OF ANY KIND EXPRESS OR IMPLIED.`)
      w.WriteHeader(code)
      fmt.Fprintf(w, "%v\n", body)
    }
  }()

  expect := func(condition bool, status int) {
    if !condition {
      panic(status)
    }
  }

  match := path_re.FindStringSubmatch(r.URL.Path)

  expect(len(match) > 0, http.StatusForbidden)

  hash := func () string {
    h := sha256.New()
    h.Write([]byte(match[1]))
    sum := h.Sum(nil)
    return hex.EncodeToString(sum)
  }()
  
  k    := datastore.NewKey(c, "Entity", hash, 0, nil)
  q    := datastore.NewQuery("Entity").KeysOnly().Filter("__key__ =", k)
  n, e := q.Count(c)

  expect(e == nil, http.StatusInternalServerError)

  stamped := n > 0
  entity  := Entity{time.Now()}
  
  status := func() int {
    switch(r.Method) {
      case "GET":
        expect(stamped, http.StatusNotFound)
        expect(datastore.Get(c, k, &entity) == nil, http.StatusInternalServerError)
        return http.StatusOK
      case "PUT":
        if stamped {
          expect(datastore.Get(c, k, &entity) == nil, http.StatusInternalServerError)
          return http.StatusOK
        } else {
          _, e := datastore.Put(c, k, &entity)
          expect(e == nil, http.StatusInternalServerError)
          return http.StatusCreated
        }
      default:
        panic(http.StatusMethodNotAllowed)
    }
  }()

  body = fmt.Sprintf("%f", float64(entity.Timestamp.UnixNano()) / 1e9)
  panic(status)
}
