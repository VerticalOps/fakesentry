package fakesentry_test

import (
	"net/http"

	"github.com/VerticalOps/fakesentry"
	raven "github.com/getsentry/raven-go"
)

func Example() {
	//Development or build flag
	var devmode bool
	//NewClient, DefaultClient
	var client *raven.Client

	if devmode {
		//Server with defaults, look at package documentation for more that can be changed.
		srv := fakesentry.NewServer()

		//Close Server on program end though it isn't needed with the default settings.
		defer srv.Close()

		//Override Sentry client Transport or make sure the underlying *http.Transport has its
		//DialX functions provided by srv.Dialer.
		client.Transport = &raven.HTTPTransport{
			Client: &http.Client{Transport: srv.Transport()},
		}

		/*
			Setting a fake DSN is optional, though it does need to be http and not https unless
			you've changed to using (or something similar).

				srv := fakesentry.NewUnstartedServer()
				go srv.ServeTLS(srv.Listener(), "cert.pem", "key.pem")

			Either way, if you don't change the underlying net.Listener completely, the Server does
			not actually use the network stack. Take a look at fakesentry docs/tests as well as fakesentry/ipc
			for more information.
		*/
		client.SetDSN(`http://thisis:myfakeauth@localhost/1`)
	}

	//All sentry/raven requests now go to the internal server.
	//fakesentry also provides a Handler type that can be used on its own, independent of the Server.
}
