/* Inspired by the tutorial http://thenewstack.io/make-a-restful-json-api-go/ logger */
package logger

import (
	"log"
	"net/http"
	"time"
)

// Logger is a handler function for http requests
// it records the request's method, the URI, and
// the total time required to execute the request.
func Logger(inner func(http.ResponseWriter, *http.Request), routeName string) func(http.ResponseWriter, *http.Request) {

	handler := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner(w, r)

		log.Printf("%s\t%s\t%s\t%s\t",
			r.Method,
			r.RequestURI,
			routeName,
			time.Since(start),
		)
	}

	return handler
}
