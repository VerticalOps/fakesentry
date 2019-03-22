package fakesentry_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/VerticalOps/fakesentry"
	raven "github.com/getsentry/raven-go"
)

func TestBasicUsage(t *testing.T) {
	eh := fakesentry.WithErrorHandler(func(w http.ResponseWriter, r *http.Request, err error) {
		t.Errorf("Error passed to handler: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	})

	const errMsg = "A bad error"
	type errJSON struct {
		Message string `json:"message"`
		EventID string `json:"event_id"`
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

	raven.DefaultClient.Transport = &raven.HTTPTransport{
		Client: &http.Client{Transport: srv.Transport()},
	}

	if err := raven.SetDSN(`http://thisis:myfakeauth@localhost/1`); err != nil {
		t.Fatalf("Unable to set raven DSN: %v", err)
	}

	eventID := raven.CaptureErrorAndWait(errors.New(errMsg), nil)
	if eventID == "" {
		t.Fatal("Did not get an event ID from raven")
	}

	if ej.Message != errMsg || ej.EventID != eventID {
		t.Fatalf("Values received from raven do not match expected values: Msg (%s) EventID (%s) JSON (%+v)", errMsg, eventID, ej)
	}
}
