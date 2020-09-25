package fakesentry_test

import (
	"github.com/VerticalOps/fakesentry"
	"github.com/getsentry/sentry-go"
)

func Example() {
	//Development or build flag
	var devmode bool
	//custom opts
	var opts sentry.ClientOptions

	if devmode {
		//Server with defaults, look at package documentation for more that can be changed.
		srv := fakesentry.NewServer()

		//Close Server on program end though it isn't needed with the default settings.
		defer srv.Close()

		//Override Sentry client Transport or make sure the underlying *http.Transport has its
		//DialX functions provided by srv.Dialer.
		opts.Dsn = `http://thisis:myfakeauth@localhost/1`
		opts.Transport = sentry.NewHTTPSyncTransport()
		opts.HTTPTransport = srv.Transport()
	}

	//All sentry requests now go to the internal server.
	//fakesentry also provides a Handler type that can be used on its own, independent of the Server.
	//look to tests for more example usage.
	if err := sentry.Init(opts); err != nil {
		//handle error
	}
}
