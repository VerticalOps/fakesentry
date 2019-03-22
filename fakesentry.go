package fakesentry

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

type Handler struct {
	logger Logger
}

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
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b, err := httputil.DumpRequest(r, false)
	if err != nil {
		h.logger.Printf("DumpRequest: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if r.ContentLength <= 0 {
		h.logger.Printf("Headers for BadRequest\n%s", b)
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if r.ContentLength > (1024 * 1024 * 32) {
		//Check for some absurd content length
		//32mb should be plenty, right?
		h.logger.Printf("Headers for BadRequest\n%s", b)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var jb []byte
	if ct := r.Header.Get("Content-Type"); ct == "application/json" {
		jb = make([]byte, r.ContentLength)

		_, err = io.ReadAtLeast(r.Body, jb, int(r.ContentLength))
		if err != nil {
			h.logger.Printf("ReadAtLeast: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else if ct == "application/octet-stream" {
		//Raven uses base64+zlib on "packets" larger than 1KB
		b64r := base64.NewDecoder(base64.StdEncoding, r.Body)

		zlr, err := zlib.NewReader(b64r)
		if err != nil {
			h.logger.Printf("zlib.NewReader: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		jb, err = ioutil.ReadAll(zlr)
		zlr.Close()
		if err != nil {
			h.logger.Printf("ReadAll: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		h.logger.Printf("Headers for BadRequest\n%s", b)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	buf := new(bytes.Buffer)
	if err = json.Indent(buf, jb, "", "  "); err != nil {
		h.logger.Printf("json.Indent: %v", err)
	}

	h.logger.Printf("\n%s%s\n", b, buf.Bytes())
}
