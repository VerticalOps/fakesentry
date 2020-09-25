package fakesentry_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/VerticalOps/fakesentry"
	"github.com/getsentry/sentry-go"
)

func TestBasicUsage(t *testing.T) {
	eh := fakesentry.WithErrorHandler(func(w http.ResponseWriter, r *http.Request, err error) {
		t.Errorf("Error passed to handler: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	})

	const errMsg = "A bad error"
	type errJSON struct {
		EventID string `json:"event_id"`
		Exe     []struct {
			Value string `json:"value"`
		} `json:"exception"`
		//just ignore everything else for now
	}

	//don't do this in production children
	ej := new(errJSON)

	mw := fakesentry.AsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jb, ok := fakesentry.FromRequest(r)
		if !ok {
			t.Fatal("Did not receive json.RawMessage as expected")
		}

		//don't do this in production children
		if err := json.Unmarshal(jb, ej); err != nil {
			t.Fatalf("json unmarshal failed: %v", err)
		}
	}))

	handler := fakesentry.NewHandler(eh, mw)

	srv := fakesentry.NewUnstartedServer()
	srv.Handler = handler

	go srv.Serve(srv.Listener())
	defer srv.Close()

	opts := sentry.ClientOptions{
		Dsn:           `http://thisis:myfakeauth@localhost/1`,
		Transport:     sentry.NewHTTPSyncTransport(),
		HTTPTransport: srv.Transport(),
	}

	if err := sentry.Init(opts); err != nil {
		t.Fatalf("Unable to init sentry: %v", err)
	}

	eventID := sentry.CaptureException(errors.New(errMsg))
	if eventID == nil {
		t.Fatal("Did not get an event ID from raven")
	}

	if len(ej.Exe) == 0 {
		t.Fatal("No captured exceptions")
	}

	if ej.Exe[0].Value != errMsg || ej.EventID != string(*eventID) {
		t.Fatalf("Values received from raven do not match expected values: Msg (%s) EventID (%s) JSON (%+v)", errMsg, *eventID, ej)
	}
}
