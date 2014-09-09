package app

import "appengine"
import "appengine/datastore"
import "fmt"
import "net/http"
import "regexp"
import "time"

var path_re = regexp.MustCompile(`^/([a-zA-Z0-9_.-]*)$`)

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
			fmt.Fprint(w, "\n", body)
		} else if body == "" {
			fmt.Fprintf(w, "%v %v\n", status.number(), status.text())
		} else {
			fmt.Fprintf(w, "%v\n", body)
		}
	}()

	ensure := func(condition bool, errorCode int) {
		if !condition {
			status = statusCode(errorCode)
			panic("ensure condition false")
		}
	}

	check := func(e error) {
		if e != nil {
			status = http.StatusInternalServerError
			panic(e)
		}
	}

	get := r.Method == "GET"

	ensure(get || r.Method == "PUT", http.StatusMethodNotAllowed)

	match := path_re.FindStringSubmatch(r.URL.Path)

	ensure(len(match) > 0, http.StatusForbidden)

	key := match[1]
	var time *time.Time

	check(datastore.RunInTransaction(c, func(c appengine.Context) error {
		time, e := getTimestamp(c, key)
		check(e)

		if get {
			ensure(time != nil, http.StatusNotFound)
			status = http.StatusOK
		} else if time == nil {
			time, e = putTimestamp(c, key)
			check(e)
			status = http.StatusCreated
		} else {
			status = http.StatusOK
		}

		return nil
	}, nil))

	body = fmt.Sprintf("%f", *time)
}
