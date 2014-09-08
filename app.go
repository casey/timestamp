package app

import "appengine"
import "appengine/datastore"
import "crypto/sha256"
import "encoding/hex"
import "fmt"
import "net/http"
import "regexp"
import "time"

var path_re = regexp.MustCompile(`^/([a-zA-Z0-9_.-]*)$`)

type Entity struct {
  Timestamp time.Time
}

func stringID(input string) string {
  sha := sha256.New()
  sha.Write([]byte(input))
  sum := sha.Sum(nil)
  return hex.EncodeToString(sum)
}

func init() {
  http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  status := statusCode(http.StatusInternalServerError)
  body := ""
  headers := make(map[string]string)
  headers["Warranty"] = `THIS IS PROVIDED "AS IS" WITHOUT WARRANTY OF ANY KIND EXPRESS OR IMPLIED.`
  headers["Content-Type"] = `text/plain; charset="utf-8"`

  defer func() {
    if e := recover(); e != nil {
      c.Errorf("handler: recovered from panic: %v", e)
    }

    for name, value := range headers {
      w.Header().Set(name, value)
    }

    w.WriteHeader(status.number())

    if status.mustNotIncludeMessageBody(r.Method) {
      fmt.Fprintf(w, "\n", body)
    } else if body == "" {
      fmt.Fprintf(w, "%v %v\n", status.number(), status.text())
    } else {
      fmt.Fprintf(w, "%v\n", body)
    }
  }()

  ensure := func(condition bool, errorCode int) {
    if !condition {
      status = statusCode(errorCode)
      panic("ensure failed")
    }
  }

  match := path_re.FindStringSubmatch(r.URL.Path)

  ensure(len(match) > 0, http.StatusForbidden)

  key := datastore.NewKey(c, "Entity", stringID(match[1]), 0, nil)
  query := datastore.NewQuery("Entity").KeysOnly().Filter("__key__ =", key)
  count, e := query.Count(c)

  ensure(e == nil, http.StatusInternalServerError)

  stamped := count > 0

  entity := Entity{time.Now()}

  switch r.Method {
  case "GET":
    ensure(stamped, http.StatusNotFound)
    ensure(datastore.Get(c, key, &entity) == nil, http.StatusInternalServerError)
    status = http.StatusOK
  case "PUT":
    if stamped {
      ensure(datastore.Get(c, key, &entity) == nil, http.StatusInternalServerError)
      status = http.StatusOK
    } else {
      _, e := datastore.Put(c, key, &entity)
      ensure(e == nil, http.StatusInternalServerError)
      status = http.StatusCreated
    }
  default:
    status = http.StatusMethodNotAllowed
  }

  if status.successful() {
    body = fmt.Sprintf("%f", float64(entity.Timestamp.UnixNano())/1e9)
  }
}
