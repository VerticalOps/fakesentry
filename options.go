package fakesentry

import (
	"context"
	"encoding/json"
	"net/http"
)

//Option is a functional option to help configure a Handler.
type Option func(*Handler)

//Logger is the basic logging interface used by this package.
type Logger interface {
	Printf(format string, arg ...interface{})
}

//WithLogger returns an Option that sets a Logger on a Handler.
func WithLogger(logger Logger) Option {
	return func(h *Handler) {
		if logger != nil {
			h.logger = logger
		}
	}
}

//ErrorHandler is called when an error occurs within Handler.ServeHTTP.
//The error is guaranteed to be non-nil, the handler may use the passed http types
//to communicate to the client what happened.
type ErrorHandler func(http.ResponseWriter, *http.Request, error)

//WithErrorHandler returns an Option for an ErrorHandler
func WithErrorHandler(eh ErrorHandler) Option {
	return func(h *Handler) {
		if eh != nil {
			h.eh = eh
		}
	}
}

type ctxKey int

const (
	rawMessage ctxKey = iota
)

//FromContext returns a json.RawMessage that was saved if AsMiddleware
//is used with the Handler type. Look to FromRequest for more information.
func FromContext(ctx context.Context) (jb json.RawMessage, ok bool) {
	jb, ok = ctx.Value(rawMessage).(json.RawMessage)
	return
}

//FromRequest retrieves the saved json.RawMessage Handler.ServeHTTP saves
//into the request context if AsMiddleware is used. Other http.Handlers can
//use this to retrieve the value and do with it as they will.
func FromRequest(r *http.Request) (json.RawMessage, bool) {
	return FromContext(r.Context())
}

//AsMiddleware takes the next http.Handler and changes the behavior of
//Handler.ServeHTTP, in that if a http.Handler is set, the parsed Sentry JSON
//value will be saved into the http.Requests context for use by the next http.Handler.
func AsMiddleware(next http.Handler) Option {
	return func(h *Handler) {
		if next != nil {
			h.next = next
		}
	}
}
