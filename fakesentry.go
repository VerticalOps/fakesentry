/*
Package fakesentry provides various functionality to run a 'fake sentry' server internal to a Go process
off or on the network. The primary motivation of this package was to provide an easy way of viewing data
your Go program may send to Sentry during run, but in a development mode without sending data to a production
DSN or otherwise setting up an entire Sentry server instance with its dependencies. As such, this should be
guarded under a development mode or build flag unless care is taken and the user has read this packages documentation.

By default, the Server type in this package will start an internal HTTP server without using the network stack
and provides a Dialer and *http.Transport that can be configured with your Sentry client to make requests to.
The Server will simply pretty-print the JSON sent by the client to STDERR. Much of this behavior can be added to
or changed. It is recommended to look at the package example and read each types documentation for more configuration.
*/
package fakesentry

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

var (
	//ErrBadMethod is passed into an ErrorHandler if the client uses an unexpected HTTP method.
	ErrBadMethod = errors.New("fakesentry: Bad HTTP method from client")

	//ErrBadContentLength is passed into an ErrorHandler if the client sends an invalid content length.
	ErrBadContentLength = errors.New("fakesentry: Bad ContentLength from client")

	//ErrBadContentType is passed into an ErrorHandler if the clients sends an unexpected content type.
	ErrBadContentType = errors.New("fakesentry: Bad ContentType from client")
)

//Handler implements a basic ServeHTTP method that can be used to accept and parse HTTP requests from
//Sentry clients. The Handler type is meant for testing, usually guarded by a development or build flag.
//It is generally unsafe to use in production unless guarded by other HTTP Handlers, such as something
//that checks the requests auth state.
type Handler struct {
	logger Logger

	eh   ErrorHandler
	next http.Handler
}

//NewHandler returns a new Handler with defaults unless Options are set.
//Unless otherwise changed, the default logger is a stdlib 'log' that prints to STDERR.
//The default ErrorHandler prints the error to the log and writes a 400 or 500 code to the client.
//If the Handler is not used in middleware then the default action is to dump the HTTP request
//as well as pretty print the Sentry JSON into the logger.
func NewHandler(opts ...Option) Handler {
	h := new(Handler)
	for _, opt := range opts {
		opt(h)
	}

	h.withDefaults()
	return *h
}

func (h *Handler) withDefaults() {
	if h.logger == nil {
		h.logger = log.New(os.Stderr, "[FAKESENTRY] ", log.LstdFlags)
	}

	if h.eh == nil {
		h.eh = h.errorHandler
	}
}

func (h Handler) errorHandler(w http.ResponseWriter, r *http.Request, err error) {
	switch err {
	case ErrBadContentLength, ErrBadContentType, ErrBadMethod:
		h.logger.Printf("%v", err)
		w.WriteHeader(http.StatusBadRequest)
	default:
		h.logger.Printf("Error from ServeHTTP: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

//ServeHTTP implements http.Handler. Please refer to the Handler type definiton
//NewHandler function for more information.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.eh(w, r, ErrBadMethod)
		return
	}

	//Check for some absurd content length
	//32mb should be plenty, right?
	if r.ContentLength <= 0 || r.ContentLength > (1024*1024*32) {
		h.eh(w, r, ErrBadContentLength)
		return
	}

	var jb []byte
	if ct := r.Header.Get("Content-Type"); ct == "application/json" {
		jb = make([]byte, r.ContentLength) //not something you want to do in production

		_, err := io.ReadAtLeast(r.Body, jb, int(r.ContentLength))
		if err != nil {
			h.eh(w, r, err)
			return
		}
	} else if ct == "application/octet-stream" {
		//Raven uses base64+zlib on "packets" larger than 1KB
		b64r := base64.NewDecoder(base64.StdEncoding, r.Body)

		zlr, err := zlib.NewReader(b64r)
		if err != nil {
			h.eh(w, r, err)
			return
		}

		jb, err = ioutil.ReadAll(zlr)
		zlr.Close()
		if err != nil {
			h.eh(w, r, err)
			return
		}
	} else {
		h.eh(w, r, ErrBadContentType)
		return
	}

	if h.next != nil {
		ctx := context.WithValue(r.Context(), rawMessage, json.RawMessage(jb))
		h.next.ServeHTTP(w, r.WithContext(ctx))
		return
	}

	b, err := httputil.DumpRequest(r, false)
	if err != nil {
		h.eh(w, r, err)
		return
	}

	buf := new(bytes.Buffer)
	if err = json.Indent(buf, jb, "", "  "); err != nil {
		h.eh(w, r, err)
	}

	h.logger.Printf("\n%s%s\n", b, buf.Bytes())
}
